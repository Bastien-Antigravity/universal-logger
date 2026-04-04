#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from asyncio import Queue as asyncioQueue, get_running_loop as asyncioGetRunningLoop, \
    get_event_loop as asyncioGetEventLoop, CancelledError as asyncioCancelledError
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from unilog import UniLog


##########################################################################
# Async config listener

class ConfigUpdateListener:
    """
    An asynchronous iterator for configuration updates.
    Returned by UniLog.on_config_update() when no callback is provided.
    """

    ##########################################################################
    # Initialization

    def __init__(self, parent: 'UniLog'):
        self._parent = parent
        self._queue = asyncioQueue()
        
        # Capture the active event loop to ensure thread-safe dispatching
        try:
            self._loop = asyncioGetRunningLoop()
        except RuntimeError:
            self._loop = asyncioGetEventLoop()
        
    # Internal thread-safe bridge to push data from Go into the Python event loop
    def _put(self, data):
        self._loop.call_soon_threadsafe(self._queue.put_nowait, data)


    ##########################################################################
    # Async Iterator protocol

    def __aiter__(self):
        # Notify the parent facade that we are now actively listening
        self._parent._async_listeners.add(self)
        return self

    async def __anext__(self):
        try:
            # Block until new data arrives via the thread-safe bridge
            return await self._queue.get()
        except asyncioCancelledError:
            # Automatic cleanup if the consumer task is cancelled
            self._parent._async_listeners.discard(self)
            raise


    ##########################################################################
    # Cleanup

    # Destruction guard to prevent memory leaks by unhooking from parent
    def __del__(self):
        if hasattr(self, '_parent'):
            self._parent._async_listeners.discard(self)
