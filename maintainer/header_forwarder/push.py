import asyncio
import logging

from maintainer import config, utils
from maintainer.bitcoin import bcoin_rpc
from maintainer.ethereum import contract, shared
from maintainer.relay_abi import ABI as relay_ABI

from typing import cast, List
from btcspv.types import RelayHeader

logger = logging.getLogger('root.header_forwarder.maintain')

HEADERS_PER_BATCH: int = 5


async def push_headers(
        best_block: RelayHeader,
        header_q: 'asyncio.Queue[RelayHeader]') -> None:
    '''Push headers from the queue to the relay contract'''
    count = 0

    # main relay tx loop
    while True:
        heads = await _get_headers_from_q(header_q)

        start_mod = heads[0]['height'] % 2016
        end_mod = heads[-1]['height'] % 2016

        # if we have a difficulty change first
        if start_mod == 0:
            await _add_diff_change(heads)
        # if we span a difficulty change
        elif start_mod > end_mod:
            pre_change = [h for h in heads if h['height'] % 2016 >= start_mod]
            post_change = [h for h in heads if h['height'] % 2016 < start_mod]

            # await all these to avoid weird race conditions.
            # Will also block if infura goes down
            if len(pre_change) != 0:
                await _add_headers(pre_change)
            if len(post_change) != 0:
                await _add_diff_change(post_change)
        else:  # if no difficulty change
            await _add_headers(heads)

        count += len(heads)
        if count >= HEADERS_PER_BATCH:
            new_best = heads[-1]
            await _update_best_digest(new_best)
            count = 0

        await asyncio.sleep(45)  # rate limit it hard


async def _get_headers_from_q(
        header_q: 'asyncio.Queue[RelayHeader]') -> List[RelayHeader]:
    '''Get whatever headers we have in the q (up to 5)'''
    headers: List[RelayHeader] = []
    # NB: wait until we have 5 headers from the q
    #     or until the queue fails to yield a header for 1 second
    while len(headers) < HEADERS_PER_BATCH:
        try:
            h = await asyncio.wait_for(header_q.get(), timeout=1)
            headers.append(h)
        except asyncio.TimeoutError:
            if len(headers) > 0:
                return headers
    return headers


async def _add_headers(headers: List[RelayHeader]) -> None:
    logger.info(f'\nsending {len(headers)} new headers\n'
                f'first is {utils.format_header(headers[0])}\n'
                f'last is {utils.format_header(headers[-1])}\n')
    nonce = next(shared.NONCE)
    anchor_or_none = await bcoin_rpc.get_header_by_hash(
        headers[0]['prevhash'].hex())
    anchor = cast(RelayHeader, anchor_or_none)

    headers_hex = ''.join(h['raw'].hex() for h in headers)

    tx = shared.make_call_tx(
        contract=config.get()['CONTRACT'],
        abi=relay_ABI,
        method='addHeaders',
        args=[anchor["raw"], headers_hex],
        nonce=nonce)
    asyncio.create_task(shared.sign_and_broadcast(tx))


async def _add_diff_change(headers: List[RelayHeader]) -> None:
    nonce = next(shared.NONCE)
    logger.info(f'\ndiff change {len(headers)} new headers,\n'
                f'first is {utils.format_header(headers[0])}\n'
                f'last is {utils.format_header(headers[-1])}\n')

    epoch_start = headers[0]['height'] - 2016
    old_start_or_none, old_end_or_none = await asyncio.gather(
        bcoin_rpc.get_header_by_height(epoch_start),
        bcoin_rpc.get_header_by_height(epoch_start + 2015),
    )

    # we know these casts won't fail
    old_start = cast(RelayHeader, old_start_or_none)
    old_end = cast(RelayHeader, old_end_or_none)
    logger.debug(f'old start is {old_start["hash"].hex()}')
    logger.debug(f'old end is {old_end["hash"].hex()}')

    headers_hex = ''.join(h['raw'].hex() for h in headers)

    tx = shared.make_call_tx(
        contract=config.get()['CONTRACT'],
        abi=relay_ABI,
        method='addHeadersWithRetarget',
        args=[
            old_start["raw"],
            old_end["raw"],
            headers_hex],
        nonce=nonce)

    asyncio.create_task(shared.sign_and_broadcast(tx))


# TODO: refactor this to not be shit
async def _update_best_digest(
        new_best: RelayHeader) -> None:
    '''Send an ethereum transaction that marks a new best known chain tip'''
    nonce = next(shared.NONCE)
    will_succeed = False

    while not will_succeed:
        current_best_digest = await contract.get_best_block()
        current_best = cast(
            RelayHeader,
            await bcoin_rpc.get_header_by_hash(current_best_digest))

        delta = new_best['height'] - current_best['height'] + 1

        # find the latest block in current's history that is an ancestor of new
        is_ancestor = False
        ancestor = current_best
        while True:
            is_ancestor = await contract.is_ancestor(
                ancestor['hash'],
                new_best['hash'])
            if is_ancestor:
                break
            ancestor = cast(
                RelayHeader,
                await bcoin_rpc.get_header_by_hash(ancestor['prevhash']))

        ancestor_le = ancestor['hash']

        tx = shared.make_call_tx(
            contract=config.get()['CONTRACT'],
            abi=relay_ABI,
            method='markNewHeaviest',
            args=[
                ancestor_le,
                current_best["raw"],
                new_best["raw"],
                delta],
            nonce=nonce)
        try:
            result = await shared.CONNECTION.preflight_tx(
                tx,
                sender=config.get()['ETH_ADDRESS'])
        except RuntimeError:
            await asyncio.sleep(10)
            continue
        will_succeed = bool(int(result, 16))
        if not will_succeed:
            await asyncio.sleep(10)

    logger.info(f'\nmarking new best\n'
                f'LCA is {ancestor["hash"].hex()}\n'
                f'previous best was {utils.format_header(current_best)}\n'
                f'new best is {utils.format_header(new_best)}\n')

    asyncio.create_task(shared.sign_and_broadcast(tx))
