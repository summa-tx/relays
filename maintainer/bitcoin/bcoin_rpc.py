import aiohttp
import logging

from maintainer import config

from maintainer.relay_types import BCoinTx
from btcspv.types import RelayHeader
from typing import Any, cast, Dict, List, Optional, Tuple, Union
S = aiohttp.ClientSession

SESSION = aiohttp.ClientSession(
    headers={"Connection": "close"}  # close the connection after each request
)

logger = logging.getLogger('root.summa_relay.bcoin_rpc')


async def close_connection() -> None:
    logger.info('closing http session')
    await SESSION.close()


async def unwrap_json(resp: aiohttp.ClientResponse) -> Dict[str, Any]:
    try:
        return cast(Dict[str, Any], await resp.json())
    except aiohttp.client_exceptions.ContentTypeError as e:
        logger.error('Failed to unwrap json from response. '
                     'Hint: is your bcoin api key correct?')
        raise e


async def _GET(route: str, session: S = SESSION) -> Tuple[int, Any]:
    '''Dispatch a GET request'''
    URL = config.get()['BCOIN_URL']

    logger.debug('get request {route}')
    full_route = f'{URL}/{route}'
    resp = await session.get(full_route)

    return resp.status, await resp.json()


async def _POST(
        route: str = '',
        payload: Dict[str, Any] = {},
        session: S = SESSION) -> Tuple[int, Any]:
    '''Dispatch a POST request'''
    URL = config.get()['BCOIN_URL']

    logger.debug(f'sending bcoin post request {payload["method"]}')
    resp = await session.post(f'{URL}/{route}', json=payload)
    status = resp.status
    resp_json = await unwrap_json(resp)

    result = None
    if resp_json is not None:
        logger.debug(f'got response {len(resp_json)}')
        result = resp_json['result'] if 'result' in resp_json else resp_json

    if status != 200:
        r = await resp.read()
        logger.error(f'Unexpected status {status} body {r!r}')
    return resp.status, result


async def _PUT(
        route: str,
        payload: Dict[str, Any],
        session: S = SESSION) -> Tuple[int, Any]:
    '''Dispatch a POST request'''
    URL = config.get()['BCOIN_URL']

    logger.debug(f'sending bcoin put request {payload["method"]}')

    resp = await session.put(f'{URL}/{route}', json=payload)
    status = resp.status
    resp_json = await unwrap_json(resp)

    result = None
    if resp_json is not None:
        logger.debug(f'got response {len(resp_json)}')
        result = resp_json['result'] if 'result' in resp_json else resp_json

    if status != 200:
        r = await resp.read()
        logger.error(f'Unexpected status {status} body {r!r}')

    return status, result


async def get_header_by_hash_le(
        hash: Union[str, bytes],
        session: S = SESSION) -> Optional[RelayHeader]:
    try:
        hash_hex = cast(bytes, hash)[::-1].hex()
    except AttributeError:
        hash_hex = bytes.fromhex(cast(str, hash))[::-1].hex()
    return await get_header_by_hash_be(hash_hex)


async def get_header_by_hash_be(
        hash: Union[str, bytes],
        session: S = SESSION) -> Optional[RelayHeader]:
    '''Gets a header by it's LE hash'''
    hash_hex: str

    try:
        hash_hex = cast(bytes, hash).hex()
    except AttributeError:
        hash_hex = cast(str, hash)

    logger.debug(f'retrieving info on {hash_hex}')
    payload = {
        'method': 'getblockheader',
        'params': [hash_hex, True]  # verbose
    }
    status, block_info_or_none = await _POST(payload=payload, session=session)
    if status != 200 or block_info_or_none is None:
        return None

    block_info = cast(dict, block_info_or_none)

    raw_payload = {
        'method': 'getblockheader',
        'params': [hash_hex, False]  # not verbose
    }
    status, raw = await _POST(payload=raw_payload, session=session)
    if status != 200:
        return None

    digest = bytes.fromhex(block_info['hash'])
    merkle_root = bytes.fromhex(block_info['merkleroot'])
    prevhash = bytes.fromhex(block_info['previousblockhash'])

    return RelayHeader(
        raw=bytes.fromhex(raw)[:80],
        hash=digest[::-1],
        height=block_info['height'],
        merkle_root=merkle_root[::-1],
        prevhash=prevhash[::-1])


async def _get_header_by_height(
        height: int,
        session: S = SESSION) -> Optional[Dict]:
    payload = {
        'method': 'getblockbyheight',
        'params': [height, True, False]  # verbose, no txns
    }
    status, block_info_or_none = await _POST(payload=payload, session=session)
    if status != 200 or block_info_or_none is None:
        return None
    return cast(dict, block_info_or_none)


async def get_header_by_height(
        height: int,
        session: S = SESSION) -> Optional[RelayHeader]:
    '''Gets useful information about a header'''
    logger.debug(f'retrieving info on block at height {height}')
    block_info_or_none = await _get_header_by_height(height, session)
    if block_info_or_none is None:
        return None

    block_info = cast(dict, block_info_or_none)
    return await get_header_by_hash_be(block_info['hash'])


async def get_chain_tips(session: S = SESSION) -> List[str]:
    logger.debug(f'retrieving info on block at chain tips')

    payload = {
        'method': 'getchaintips'
    }
    status, res = await _POST(payload=payload, session=session)
    if status != 200:
        raise RuntimeError(f'Unexpected status in get_chain_tips: {status}')
    return [a['hash'] for a in res]


async def get_tx(tx_id: bytes, session: S = SESSION) -> Optional[BCoinTx]:
    route = f'tx/{tx_id[::-1].hex()}'  # make BE
    status, res = await _GET(route, session)
    if status != 200:
        return None
    return cast(BCoinTx, res)
