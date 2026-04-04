#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import ctypes
import json
import inspect
from models import LogLevel
from lib_loader import lib, CALLBACK_TYPE
from listeners import ConfigUpdateListener

class UniLog:
    """
    Python Facade for the Universal Logger (Go) shared library.
    Provides integrated configuration management and high-performance logging.
    """

    def __init__(self, config_profile="standalone", app_name="python-app", 
                 logger_profile="standard", log_level="info"):
        if not lib:
            raise RuntimeError("libunilog shared library not found. Please ensure it is built.")
        
        # Convert string log level to int for Go
        level_val = LogLevel.from_str(log_level) if isinstance(log_level, str) else int(log_level)

        self._handle = lib.UniLog_Init(
            config_profile.encode('utf-8'), 
            app_name.encode('utf-8'), 
            logger_profile.encode('utf-8'), 
            ctypes.c_int(level_val)
        )
        self._callback_ref = None # Keep reference to avoid GC
        self._sync_subscribers = []
        self._async_listeners = set()
        self._initialized_bridge = False

    def close(self):
        """Manually release the logger session and associated resources."""
        if hasattr(self, '_handle') and self._handle:
            lib.UniLog_Close(self._handle)
            self._handle = None

    def __del__(self):
        self.close()

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()

    async def __aenter__(self):
        return self.__enter__()

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        return self.__exit__(exc_type, exc_val, exc_tb)

    # --- Configuration Methods ---

    def get_config(self, section: str, key: str, default: str = None) -> str:
        """Retrieve a configuration value from the distributed config service."""
        res = lib.UniLog_Config_Get(self._handle, section.encode('utf-8'), key.encode('utf-8'))
        if res is None:
            return default
        return res.decode('utf-8')

    def set_config(self, section: str, key: str, value: str):
        """Update a configuration value in the memory configuration."""
        lib.UniLog_Config_Set(self._handle, section.encode('utf-8'), key.encode('utf-8'), value.encode('utf-8'))

    def _dispatch_update(self, json_data):
        """Internal bridge called from Go shared library background thread."""
        try:
            data = json.loads(json_data.decode('utf-8'))
            
            # 1. Dispatch to synchronous subscribers
            for cb in self._sync_subscribers:
                try:
                    cb(data)
                except:
                    pass
            
            # 2. Dispatch to asynchronous listeners (thread-safe)
            for listener in list(self._async_listeners): # Copy list to avoid concurrent mutation
                listener._put(data)
        except Exception:
            pass

    def on_config_update(self, callback=None) -> ConfigUpdateListener:
        """
        Registers a mechanism for configuration updates.
        
        Args:
            callback: (Optional) A standard Python function. If provided, 
                      registers a traditional synchronous callback.
        
        Returns:
            None: If a callback was provided.
            ConfigUpdateListener: If NO callback was provided. Use with 'async for'.
        """
        # Lazy initialization of the single C bridge
        if not self._initialized_bridge:
            self._callback_ref = CALLBACK_TYPE(self._dispatch_update)
            lib.UniLog_OnMemConfUpdate(self._handle, self._callback_ref)
            self._initialized_bridge = True

        if callback is not None:
            self._sync_subscribers.append(callback)
            return None
        
        return ConfigUpdateListener(self)

    # --- Logging Methods ---

    def set_level(self, level):
        """Change the current log level dynamically."""
        if isinstance(level, str):
            level = LogLevel.from_str(level)
        lib.UniLog_SetLevel(self._handle, int(level))

    def _log(self, level: int, msg: str):
        # Capture caller information for structured logging
        caller = inspect.stack()[2]
        filename = os.path.basename(caller.filename)
        lineno = str(caller.lineno)
        function = caller.function
        module = caller.frame.f_globals.get('__name__', 'unknown')

        lib.UniLog_LogWithMetadata(
            self._handle, 
            int(level), 
            str(msg).encode('utf-8'), 
            filename.encode('utf-8'), 
            lineno.encode('utf-8'), 
            function.encode('utf-8'), 
            module.encode('utf-8')
        )

    def debug(self, msg): self._log(LogLevel.DEBUG, msg)
    def info(self, msg): self._log(LogLevel.INFO, msg)
    def warning(self, msg): self._log(LogLevel.WARNING, msg)
    def error(self, msg): self._log(LogLevel.ERROR, msg)
    def critical(self, msg): self._log(LogLevel.CRITICAL, msg)
