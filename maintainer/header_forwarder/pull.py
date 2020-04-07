import asyncio
import logging

from relay.bitcoin import bcoin_rpc

from typing import cast
from btcspv.types import RelayHeader

logger = logging.getLogger('root.header_forwarder.poll')


async def pull_headers(
        latest: RelayHeader,
        q: 'asyncio.Queue[RelayHeader]') -> None:
    latest_height = latest['height'] + 1
    last_added: RelayHeader = cast(RelayHeader, {})

    while True:
        header_or_none = await bcoin_rpc.get_header_by_height(latest_height)

        # sleep then loop again if no header received
        if header_or_none is None:
            logger.debug('sleeping at tip')
            await asyncio.sleep(60)
            continue

        header = cast(RelayHeader, header_or_none)
        if header != last_added:
            await q.put(header)  # will block if q is full
            last_added = header
            latest_height = latest_height + 1
