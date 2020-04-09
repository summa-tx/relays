import asyncio
import logging

from ether import calldata, ethrpc

from maintainer import config

from ether.ether_types import Receipt
from ether.transactions import UnsignedEthTx
from typing import Any, cast, Dict, Iterator, List, Optional

logger = logging.getLogger('root.summa_relay.shared_eth')


GWEI = 1000000000
DEFAULT_GAS = 500_000
DEFAULT_GAS_PRICE = 2 * GWEI

CONNECTION: ethrpc.BaseRPC
NONCE: Iterator[int]  # yields ints, takes no sends


def _nonce(i: int) -> Iterator[int]:
    '''Infinite generator for nonces'''
    index = i
    while 1:
        yield index
        index += 1


async def init() -> None:
    '''Set up a connection to the interwebs'''
    global CONNECTION

    c = config.get()
    network = c['NETWORK']
    project_id = c['PROJECT_ID']
    uri = c['ETHER_URL']
    force_https = project_id != ''

    logger.info(f'contract is {c["CONTRACT"]}')

    CONNECTION = ethrpc.get_client(
        network=network,
        infura_key=project_id,
        uri=uri,
        logger=logger.getChild('ethrpc'),
        force_https=force_https)

    await CONNECTION.open()

    if c['PRIVKEY'] is None and c['GETH_UNLOCK'] is None:
        logger.warn(
            'No ethereum privkey found in env config. Txns will error')
    else:
        global NONCE
        address = cast(str, c['ETH_ADDRESS'])
        n = await CONNECTION.get_nonce(address)
        NONCE = _nonce(n)
        logger.info(f'nonce is {n}')


async def close_connection() -> None:
    try:
        global CONNECTION
        await CONNECTION.close()
    except NameError:
        pass


async def sign_and_broadcast(
        tx: UnsignedEthTx,
        ignore_result: bool = False) -> None:
    '''Sign an ethereum transaction and broadcast it to the network'''
    c = config.get()
    privkey = c['PRIVKEY']
    address = c['ETH_ADDRESS']
    unlock_code = c['GETH_UNLOCK']

    if privkey is None and unlock_code is None:
        raise RuntimeError('Attempted to sign tx without access to key')

    if privkey is None:
        logger.debug('signing with ether node')
        await CONNECTION._RPC(
            'personal_unlockAccount',
            [address, unlock_code])
        tx_id = await CONNECTION.send_transaction(cast(str, address), tx)
    else:
        logger.debug('signing with local key')
        signed = tx.sign(cast(bytes, privkey))
        serialized = signed.serialize_hex()
        tx_id = await CONNECTION.broadcast(serialized)

    logger.info(f'dispatched transaction {tx_id}')
    if not ignore_result:
        asyncio.ensure_future(_track_tx_result(tx_id))


def make_call_tx(
        contract: str,
        abi: List[Dict[str, Any]],
        method: str,
        args: List[Any],
        nonce: int,
        value: int = 0,
        gas: int = DEFAULT_GAS,
        gas_price: int = DEFAULT_GAS_PRICE) -> UnsignedEthTx:
    '''
    Sends tokens to a recipient
    Args:
        contract      (str): address of contract being called
        abi          (dict): contract ABI
        method        (str): the name of the method to call
        args         (list): the arguments to the method call
        nonce         (int): the account nonce for the txn
        value         (int): ether in wei
        gas_price     (int): the price of gas in wei or gwei
    Returns:
        (UnsignedEthTx): the unsigned tx object
    '''
    logger.debug(f'making tx call {method} on {contract} '
                 f'with value {value} and {len(args)} args')

    gas_price = _adjust_gas_price(gas_price)
    chainId = config.get()['CHAIN_ID']

    data = calldata.call(
        method,
        args,
        abi)

    txn = UnsignedEthTx(
        to=contract,
        value=value,
        gas=gas,
        gasPrice=gas_price,
        nonce=nonce,
        data=data,
        chainId=chainId)

    return txn


def _adjust_gas_price(gas_price: int) -> int:
    '''
    We accept gas price in GWEI or in WEI.
    This adjusts, and ensures we error if it's high.
    Args:
        gas_price (int): the user-provided gas price
    Returns:
        (int): the adjusted price
    '''
    if gas_price < GWEI:
        gas_price = gas_price * GWEI
    if gas_price > 1000 * GWEI:
        logger.error('rejecting high gas price')
        raise ValueError(
            'very high gas price detected: {} gwei'.format(gas_price / GWEI))
    return gas_price


async def _track_tx_result(tx_id: str) -> None:
    '''Keep track of the result of a transaction by polling every 25 seconds'''
    receipt_or_none: Optional[Receipt] = None

    for _ in range(20):
        await asyncio.sleep(30)
        receipt_or_none = await CONNECTION.get_tx_receipt(tx_id)
        if receipt_or_none is not None:
            break

    if receipt_or_none is None:
        raise RuntimeError(f'No receipt after 10 minutes: {tx_id}')

    receipt = cast(Receipt, receipt_or_none)
    logger.info(f'Receipt for {tx_id} status is {receipt["status"]}')

    if receipt['status'] != '0x1':
        raise RuntimeError(f'Failed tx: {receipt["transactionHash"]}')
