#!/usr/bin/env python
# coding:utf-8
"""
=============================================================================
ESSENTIAL PROCESS:
Asynchronous event listeners and subscription streams for configuration updates
and telemetry notifications in the Python UniLog client.

DATA FLOW:
1. Input: Raw CGO callback events from Go background workers.
2. Logic: Buffers updates into thread-safe queues and yields items to async generators.
3. Output: Async iterable stream of configuration and notification payloads.

KEY PARAMETERS:
- queue_size: Maximum capacity of the background update event buffer.
=============================================================================
"""

from asyncio import CancelledError as asyncioCancelledError, Queue as asyncioQueue, get_event_loop as asyncioGetEventLoop, get_running_loop as asyncioGetRunningLoop
from typing import Any as typeAny, Dict as typeDict, TYPE_CHECKING as typeChecking

if typeChecking:
    from unilog.facade import UniLog

# -----------------------------------------------------------------------------------------------


class ConfigUpdateListener:
    """
    An asynchronous iterator for configuration updates.
    Returned by UniLog.on_config_update() when no callback is provided.
    """

    Name = "ConfigUpdateListener"

    # -----------------------------------------------------------------------------------------------

    def __init__(self, parent: 'UniLog') -> None:
        self._parent: 'UniLog' = parent
        self._queue: asyncioQueue[typeDict[str, typeAny]] = asyncioQueue()

        # Capture the active event loop to ensure thread-safe dispatching
        try:
            self._loop = asyncioGetRunningLoop()
        except RuntimeError:
            self._loop = asyncioGetEventLoop()

    # -----------------------------------------------------------------------------------------------

    # Internal thread-safe bridge to push data from Go into the Python event loop
    def _put(self, data: typeDict[str, typeAny]) -> None:
        self._loop.call_soon_threadsafe(self._queue.put_nowait, data)

    # -----------------------------------------------------------------------------------------------

    def __aiter__(self) -> 'ConfigUpdateListener':
        # Notify the parent facade that we are now actively listening
        self._parent._async_listeners.add(self)
        return self

    # -----------------------------------------------------------------------------------------------

    async def __anext__(self) -> typeDict[str, typeAny]:
        try:
            # Block until new data arrives via the thread-safe bridge
            return await self._queue.get()
        except asyncioCancelledError:
            # Automatic cleanup if the consumer task is cancelled
            self._parent._async_listeners.discard(self)
            raise

    # -----------------------------------------------------------------------------------------------

    # Destruction guard to prevent memory leaks by unhooking from parent
    def __del__(self) -> None:
        if hasattr(self, '_parent'):
            self._parent._async_listeners.discard(self)
