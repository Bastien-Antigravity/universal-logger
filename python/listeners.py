#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import asyncio
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from unilog import UniLog

class ConfigUpdateListener:
    """
    An asynchronous iterator for configuration updates.
    Returned by UniLog.on_config_update() when no callback is provided.
    """
    def __init__(self, parent: 'UniLog'):
        self._parent = parent
        self._queue = asyncio.Queue()
        # Capture the loop where this listener is created
        try:
            self._loop = asyncio.get_running_loop()
        except RuntimeError:
            self._loop = asyncio.get_event_loop()
        
    def _put(self, data):
        """Internal thread-safe producer."""
        self._loop.call_soon_threadsafe(self._queue.put_nowait, data)

    def __aiter__(self):
        # Register this iterator instance with the parent dispatcher
        self._parent._async_listeners.add(self)
        return self

    async def __anext__(self):
        try:
            return await self._queue.get()
        except asyncio.CancelledError:
            # Unregister on cancellation
            self._parent._async_listeners.discard(self)
            raise

    def __del__(self):
        # Ensure we discard from parent on cleanup
        if hasattr(self, '_parent'):
            self._parent._async_listeners.discard(self)
