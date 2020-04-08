import logging
from ether import abi, calldata, events

from maintainer import config
from maintainer.ethereum import shared
from maintainer.relay_abi import ABI as relay_ABI

EXPIRED = events._make_topic0(
    abi.find('SubscriptionExpired', relay_ABI)[0])
CLOSED = events._make_topic0(
    abi.find('RequestClosed', relay_ABI)[0])
FILLED = events._make_topic0(
    abi.find('RequestFilled', relay_ABI)[0])

logger = logging.getLogger('root.summa_relay.eth_contract')


async def find_height(digest_le: bytes) -> int:
    data = calldata.call(
        "findHeight",
        [digest_le],
        relay_ABI)
    res = await shared.CONNECTION._RPC(
        method='eth_call',
        params=[
            {
                'from': config.get()['ETH_ADDRESS'],
                'to': config.get()['CONTRACT'],
                'data': f'0x{data.hex()}'
            },
            'latest'  # block height parameter
        ]
    )
    logger.debug(f'findHeight for {digest_le.hex()} is {res}')
    return int(res, 16)


async def has_block(digest_le: bytes) -> bool:
    '''Check if the relay knows of a block'''
    height = await find_height(digest_le)
    logger.debug(f'height is {height}')
    return height != 0


async def is_ancestor(
        ancestor: bytes,
        descendant: bytes,
        limit: int = 240) -> bool:
    '''Determine if ancestor precedes descendant'''
    data = calldata.call(
        "isAncestor",
        [ancestor, descendant, limit],
        relay_ABI)
    res = await shared.CONNECTION._RPC(
        method='eth_call',
        params=[
            {
                'from': config.get()['ETH_ADDRESS'],
                'to': config.get()['CONTRACT'],
                'data': f'0x{data.hex()}'
            },
            'latest'  # block height parameter
        ]
    )
    # returned as 0x-prepended hex string representing 32 bytes
    return bool(int(res, 16))


async def get_best_block() -> str:
    '''
    Get the contract's marked best known digest.
    Counterintuitively, the contract may know of a better digest
      that hasn't been marked yet

    returns LE digest
    '''
    f = abi.find('getBestKnownDigest', relay_ABI)[0]
    selector = calldata.make_selector(f)
    res = await shared.CONNECTION._RPC(
        method='eth_call',
        params=[
            {
                'from': config.get()['ETH_ADDRESS'],
                'to': config.get()['CONTRACT'],
                'data': f'0x{selector.hex()}'
            },
            'latest'  # block height parameter
        ]
    )
    return res[2:]  # block-explorer format
