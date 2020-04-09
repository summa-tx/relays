from riemann import tx
from btcspv.types import RelayHeader
from maintainer.relay_types import RelayRequest


def format_header(h: RelayHeader) -> str:
    return f'height {h["height"]} hash {h["hash"].hex()}'


def format_request(r: RelayRequest) -> str:
    return (
        f'id {r["request_id"]}\n'
        f'pays {r["pays_addr"]}\n'
        f'spends {r["spends"].hex()}')


def extract_vin(t: tx.Tx) -> bytes:
    '''Get the length-prefixed input vector from a tx'''
    b = bytearray([len(t.tx_ins)])
    for tx_in in t.tx_ins:
        b.extend(tx_in)
    return bytes(b)


def extract_vout(t: tx.Tx) -> bytes:
    '''Get the length-prefixed output vector from a tx'''
    b = bytearray([len(t.tx_outs)])
    for tx_out in t.tx_outs:
        b.extend(tx_out)
    return bytes(b)


def reverse_hex_bytes(hex_bytes: str) -> str:
    '''Take a hex string, return it with the bytes in opposite order'''
    return bytes.fromhex(hex_bytes)[::-1].hex()
