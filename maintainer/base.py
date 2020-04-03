import sys
import signal
import asyncio
import logging
from pathlib import Path
from functools import partial
from dotenv import load_dotenv

from relay import config

from typing import Awaitable, Callable
from asyncio.events import AbstractEventLoop

AsyncFunction = Callable[[], Awaitable[None]]


def registerFileHandler(name: str, logger: logging.Logger) -> None:
    if sys.platform.startswith('win'):
        raise NotImplementedError('Windows not supported')  # pragma: nocover
    logDir = Path.home() / '.summa' / 'relays'
    logDir.mkdir(parents=True, exist_ok=True)

    logPath = logDir / name

    formatter = logging.Formatter(
        fmt='%(asctime)6s %(name)s: %(levelname)s %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S')
    fh = logging.FileHandler(logPath)
    fh.setFormatter(formatter)
    fh.setLevel(logging.DEBUG)
    logger.addHandler(fh)


def get_env_name(default: str) -> str:
    '''Checks for a argv-passed env name, formats the default otherwise'''
    if len(sys.argv) > 1:
        return sys.argv[1]
    return f'.{default}.env'


def set_config(env_name: str) -> None:
    '''Load dotfiles and set the config object'''
    # Load config from .env file(s)
    # Load a base .env first, then override with an app-specified version
    path = Path(__file__).parent / 'config'
    base_env = path / '.env'
    load_dotenv(base_env, override=True)
    load_dotenv(path / env_name, override=True)
    config.set()


def main(
        run: AsyncFunction,
        teardown: AsyncFunction,
        name: str,
        logger: logging.Logger) -> None:
    '''Template for small, headless, async applications'''
    logger.info(f'Setting config {name}')
    env_name = get_env_name(default=name)
    set_config(env_name=env_name)
    registerFileHandler(name=env_name, logger=logger)
    logger.info('Starting relay')

    loop = asyncio.get_event_loop()

    # set up graceful exit
    signals = (signal.SIGHUP, signal.SIGTERM, signal.SIGINT)
    for s in signals:
        loop.add_signal_handler(
            s, lambda s=s: asyncio.create_task(
                shutdown(loop, logger, teardown, signal=s)))

    handler = partial(handle_exception, logger=logger, teardown=teardown)
    loop.set_exception_handler(handler)

    asyncio.ensure_future(run())
    loop.run_forever()


async def shutdown(  # type: ignore[no-untyped-def]
        loop: AbstractEventLoop,
        logger: logging.Logger,
        teardown: AsyncFunction,
        signal=None  # a named enum of ints
) -> None:
    '''Cancel active tasks for shutdown'''
    if signal:
        logger.info(f'Received exit signal {signal.name}')
    else:
        logger.info('Unexpeced shutdown initiated')
        await asyncio.sleep(5)  # stall error loops

    if teardown:
        try:
            await teardown()
        except Exception:
            logger.exception('Error during teardown function')
            logger.error('Exiting uncleanly')
            sys.exit(1)

    tasks = [t for t in asyncio.Task.all_tasks() if t is not
             asyncio.current_task()]

    logger.info(f'Cancelling {len(tasks)} tasks')
    [task.cancel() for task in tasks]

    try:
        await asyncio.gather(*tasks, return_exceptions=True)
    except Exception:
        logger.exception('Error during loop task cancellation')
        logger.error('Exiting uncleanly')
        sys.exit(1)

    loop.stop()


def handle_exception(  # type: ignore[no-untyped-def]
        loop: AbstractEventLoop,
        context,  # don't worry about it
        logger: logging.Logger,
        teardown: AsyncFunction):
    '''Global exception handler. Gets all unhandled exceptions from tasks'''
    # context['message'] will always be there; but context['exception'] may not
    if 'exception' in context:
        # reraise so we can log it
        try:
            raise context['exception']
        except Exception:
            logger.exception('Caught exception')
    else:
        logger.error(f'Caught exception: {context["message"]}')

    logger.info('Shutting down')
    asyncio.create_task(shutdown(loop, logger, teardown))
