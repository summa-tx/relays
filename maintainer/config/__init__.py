import os

from ether import crypto

from typing import cast, Tuple, Optional
from relay.relay_types import RelayConfig

CONFIG: RelayConfig


def is_infura() -> bool:
    pid = get()['PROJECT_ID']
    return not pid == ''


def _set_keys() -> Tuple[Optional[bytes], Optional[bytes], Optional[str]]:
    # Keys
    PRIVKEY: Optional[bytes]
    PUBKEY: Optional[bytes]
    ETH_ADDRESS: Optional[str]

    PRIVKEY_HEX = os.environ.get('SUMMA_RELAY_OPERATOR_KEY', None)
    try:
        PRIVKEY = bytes.fromhex(cast(str, PRIVKEY_HEX))
    except (ValueError, TypeError):  # hex errors or is None
        PRIVKEY = None

    PUBKEY = crypto.priv_to_pub(PRIVKEY) if PRIVKEY else None

    if PRIVKEY:
        ETH_ADDRESS = crypto.priv_to_addr(PRIVKEY)
    else:
        ETH_ADDRESS = os.environ.get('OPERATOR_ADDRESS', None)

    return PRIVKEY, PUBKEY, ETH_ADDRESS


def _set_net() -> Tuple[str, int]:
    CHAIN_IDS = {'mainnet': 1, 'ropsten': 3, 'kovan': 42}
    NETWORK = os.environ.get('SUMMA_RELAY_ETH_NETWORK', 'ropsten')
    if NETWORK in CHAIN_IDS:
        CHAIN_ID = CHAIN_IDS[NETWORK]
    else:
        try:
            CHAIN_ID = int(
                os.environ.get('SUMMA_RELAY_ETH_CHAIN_ID'))  # type: ignore
        except (ValueError, TypeError):
            CHAIN_ID = 1

    return NETWORK, CHAIN_ID


def get() -> RelayConfig:
    return CONFIG


def set() -> RelayConfig:
    BCOIN_HOST = os.environ.get('SUMMA_RELAY_BCOIN_HOST', '127.0.0.1')
    API_KEY = os.environ.get('SUMMA_RELAY_BCOIN_API_KEY', '')
    BCOIN_PORT = os.environ.get('SUMMA_RELAY_BCOIN_PORT', '8332')

    ETHER_HOST = os.environ.get('SUMMA_RELAY_ETHER_HOST', '127.0.0.1')
    ETHER_PORT = os.environ.get('SUMMA_RELAY_ETHER_PORT', '8545')

    GETH_UNLOCK = os.environ.get('SUMMA_RELAY_GETH_UNLOCK', None)

    PRIVKEY, PUBKEY, ETH_ADDRESS = _set_keys()

    NETWORK, CHAIN_ID = _set_net()

    global CONFIG
    CONFIG = RelayConfig(
        PRIVKEY=PRIVKEY,
        PUBKEY=PUBKEY,
        ETH_ADDRESS=ETH_ADDRESS,
        NETWORK=NETWORK,
        CHAIN_ID=CHAIN_ID,
        API_KEY=API_KEY,
        BCOIN_HOST=BCOIN_HOST,
        BCOIN_PORT=BCOIN_PORT,
        ETHER_HOST=ETHER_HOST,
        ETHER_PORT=ETHER_PORT,
        GETH_UNLOCK=GETH_UNLOCK,
        ETHER_URL=f'http://{ETHER_HOST}:{ETHER_PORT}',
        BCOIN_URL=f'http://x:{API_KEY}@{BCOIN_HOST}:{BCOIN_PORT}',
        BCOIN_WS_URL=f'ws://{BCOIN_HOST}:{BCOIN_PORT}',
        PROJECT_ID=os.environ.get('SUMMA_RELAY_INFURA_KEY', ''),
        CONTRACT=os.environ.get('SUMMA_RELAY_CONTRACT', ''),
    )

    return CONFIG
