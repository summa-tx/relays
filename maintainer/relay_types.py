# from riemann import tx

from mypy_extensions import TypedDict
from typing import Any, Dict, List, Optional


class BcoinOutpoint(TypedDict):
    hash: str
    index: int


class RelayRequest(TypedDict):
    request_id: int
    pays_addr: str
    pays_value: int
    spends: bytes
    pays: bytes  # length prefixed


class BcoinRequest(TypedDict):
    id: str
    address: str
    value: int
    spends: BcoinOutpoint
    pays: str  # not length prefixed


class Notification(TypedDict):
    height: int


class RelayNotification(Notification):
    tx_id: bytes
    satisfied: List[int]


class BcoinNotification(Notification):
    txid: str
    satisfied: List[str]


class BCoinTx(TypedDict):
    hash: str
    witnessHash: str
    fee: int
    rate: int
    mtime: int
    height: int
    block: Optional[str]
    time: int
    index: int
    version: int
    inputs: List[Dict[str, Any]]
    outputs: List[Dict[str, Any]]
    locktime: int
    hex: str
    confirmations: int


class RelayConfig(TypedDict):
    PRIVKEY: Optional[bytes]
    PUBKEY: Optional[bytes]
    ETH_ADDRESS: Optional[str]
    NETWORK: str
    CHAIN_ID: int
    API_KEY: str
    BCOIN_HOST: str
    BCOIN_PORT: str
    ETHER_HOST: str
    ETHER_PORT: str
    ETHER_URL: str
    GETH_UNLOCK: Optional[str]
    BCOIN_URL: str
    BCOIN_WS_URL: str
    PLUGIN_WS_URL: str
    PROJECT_ID: str
    CONTRACT: str
