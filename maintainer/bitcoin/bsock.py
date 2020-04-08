import asyncio
import logging
import socketio

from maintainer import config

from typing import Dict, List, Tuple, Union

logger = logging.getLogger('root.summa_relay.bsock')

cl = logger.getChild('async_client')
sio = socketio.AsyncClient(logger=cl)

SioData = Union[str, bytes, Dict, List, Tuple]


async def get_connection() -> socketio.AsyncClient:
    logger.info('opening bsock ws session')
    if not sio.connected:
        await sio.connect(config.get()['BCOIN_WS_URL'], transports='websocket')
    return sio


async def close_connection() -> None:
    logger.info('closing bsock ws session')
    if sio.connected:
        await sio.disconnect()


@sio.event
async def connect() -> None:
    await sio.call('auth', config.get()['API_KEY'])
    logger.info(f'connected and authed')


@sio.event
async def disconnect() -> None:
    logger.info('bsock disconnected')


@sio.event
async def tx(data: SioData) -> None:
    logger.info('message\t', data)


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    asyncio.ensure_future(get_connection())
    asyncio.get_event_loop().run_forever()
