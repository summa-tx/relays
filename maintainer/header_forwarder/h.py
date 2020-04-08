import sys
import asyncio
import logging

from maintainer import base, utils
from maintainer.bitcoin import bcoin_rpc, bsock
from maintainer.ethereum import contract, shared
from maintainer.header_forwarder import pull, push

from typing import cast
from btcspv.types import RelayHeader

logger = logging.getLogger('root.header_forwarder')
logging.basicConfig(
    format='%(asctime)6s %(name)s: %(levelname)s %(message)s',
    level=logging.INFO,
    datefmt='%Y-%m-%d %H:%M:%S')


async def run() -> None:
    header_q: 'asyncio.Queue[RelayHeader]' = asyncio.Queue(maxsize=50)
    await shared.init()

    latest_digest = await contract.get_best_block()

    if len(latest_digest) != 64:
        raise ValueError(
            'Expected 32 byte digest from contract. '
            f'Received {len(latest_digest) // 2} bytes instead. '
            'Hint: is this account authorized?')

    latest_or_none = await bcoin_rpc.get_header_by_hash(latest_digest)
    if latest_or_none is None:
        raise ValueError(
            'Relay\'s latest digest is not known to the Bitcoin node. '
            f'Got {latest_digest}. '
            'Hint: is your node on the same Bitcoin network as the relay?')
    latest = cast(RelayHeader, latest_or_none)
    better_or_same = cast(
        RelayHeader,
        await bcoin_rpc.get_header_by_height(latest['height']))

    # see if there's a better block at that height
    # if so, crawl backwards
    while latest != better_or_same:
        latest = cast(
            RelayHeader,
            await bcoin_rpc.get_header_by_hash(latest['prevhash']))
        better_or_same = cast(
            RelayHeader,
            await bcoin_rpc.get_header_by_height(latest['height']))

    logger.info(
        f'latest is {utils.format_header(latest)}')

    asyncio.create_task(pull.pull_headers(latest, header_q))
    asyncio.create_task(push.push_headers(latest, header_q))


async def teardown() -> None:
    coros = [
        # close http session
        bcoin_rpc.close_connection(),
        # close socketio connection
        bsock.close_connection(),
        # close infura websocket
        shared.close_connection()
    ]

    await asyncio.gather(*coros, return_exceptions=True)


if __name__ == '__main__':
    try:
        name = 'header_forwarder'
        base.main(run=run, logger=logger, name=name, teardown=teardown)
    except Exception:
        logger.exception('---- Fatal Exception ----')
        sys.exit(1)
