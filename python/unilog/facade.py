#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os.path import basename as osPathBasename
from ctypes import c_int as ctypeC_int
from inspect import stack as inspectStack
from json import loads as jsonLoads
from asyncio import get_running_loop as asyncioGetRunningLoop

from .models import LogLevel
from .listeners import ConfigUpdateListener
from .lib_loader import lib, CALLBACK_TYPE


class UniLog:
    """
    Python Facade for the Universal Logger (Go) shared library.
    Provides integrated configuration management and high-performance logging.
    """

    def __init__(self, config_profile="standalone", app_name="python-app", 
                 logger_profile="standard", log_level="info", use_local_notifier=False):
        if not lib:
            raise RuntimeError("libunilog shared library not found. Please ensure it is built.")
        
        # Convert string log level to int for Go
        level_val = LogLevel.from_str(log_level) if isinstance(log_level, str) else int(log_level)

        self._handle = lib.UniLog_Init(
            config_profile.encode('utf-8'), 
            app_name.encode('utf-8'), 
            logger_profile.encode('utf-8'), 
            ctypeC_int(level_val),
            ctypeC_int(1 if use_local_notifier else 0)
        )
        self._callback_ref = None # Keep reference to avoid GC
        self._sync_subscribers = []
        self._async_listeners = set()
        self._initialized_bridge = False


    ##########################################################################
    # Logging Methods
    
    def debug(self, msg): self._log(LogLevel.DEBUG, msg)
    def info(self, msg): self._log(LogLevel.INFO, msg)
    def warning(self, msg): self._log(LogLevel.WARNING, msg)
    def error(self, msg): self._log(LogLevel.ERROR, msg)
    def critical(self, msg): self._log(LogLevel.CRITICAL, msg)

    async def async_debug(self, msg): await self._async_log(LogLevel.DEBUG, msg)
    async def async_info(self, msg): await self._async_log(LogLevel.INFO, msg)
    async def async_warning(self, msg): await self._async_log(LogLevel.WARNING, msg)
    async def async_error(self, msg): await self._async_log(LogLevel.ERROR, msg)
    async def async_critical(self, msg): await self._async_log(LogLevel.CRITICAL, msg)

    # Log level setter
    def set_level(self, level):
        """Change the current log level dynamically."""
        if isinstance(level, str):
            level = LogLevel.from_str(level)
        lib.UniLog_SetLevel(self._handle, int(level))
        

    ##########################################################################
    # Config Methods ---

    def get_config(self, section: str, key: str, default: str = None) -> str:
        """Retrieve a configuration value from the distributed config service."""
        res = lib.UniLog_Config_Get(self._handle, section.encode('utf-8'), key.encode('utf-8'))
        if res is None:
            return default
        return res.decode('utf-8')

    def set_config(self, section: str, key: str, value: str):
        """Update a configuration value in the memory configuration."""
        lib.UniLog_Config_Set(self._handle, section.encode('utf-8'), key.encode('utf-8'), value.encode('utf-8'))

    # Trigger on_config_update regarding the caller and caller method
    def _dispatch_update(self, json_data):
        """Internal bridge called from Go shared library background thread."""
        print(f"!!! _dispatch_update entered with: {json_data}")
        try:
            raw_val = json_data.decode('utf-8')
            print(f"!!! Decoding successful: {raw_val}")
            data = jsonLoads(raw_val)
            
            # 1. Dispatch to synchronous subscribers
            for cb in self._sync_subscribers:
                try:
                    print(f"!!! Calling sync subscriber: {cb}")
                    cb(data)
                except Exception as e:
                    print(f"!!! Sync subscriber error: {e}")
            
            # 2. Dispatch to asynchronous listeners (thread-safe)
            for listener in list(self._async_listeners): # Copy list to avoid concurrent mutation
                print(f"!!! Calling async listener: {listener}")
                listener._put(data)
        except Exception as e:
            print(f"!!! _dispatch_update EXCEPTION: {e}")


    ##########################################################################
    # Trigger on_config_update

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
            print(f"!!! Python: Registering C bridge for handle: {self._handle}")
            self._callback_ref = CALLBACK_TYPE(self._dispatch_update)
            print(f"!!! Python: callback_ref created: {self._callback_ref}")
            lib.UniLog_OnMemConfUpdate(self._handle, self._callback_ref)
            print("!!! Python: UniLog_OnMemConfUpdate call finished.")
            self._initialized_bridge = True

        if callback is not None:
            self._sync_subscribers.append(callback)
            return None
        
        return ConfigUpdateListener(self)


    ##########################################################################
    # Local Notifier Methods

    def on_notification(self, callback):
        """
        Registers a callback for local notifications.
        The callback will receive a dictionary parsed from the notification JSON.
        """
        def _bridge_cb(json_data):
            try:
                data = jsonLoads(json_data.decode('utf-8'))
                callback(data)
            except Exception as e:
                print(f"!!! on_notification EXCEPTION: {e}")

        # Keep a reference to the bridge callback to avoid GC
        self._notif_callback_ref = CALLBACK_TYPE(_bridge_cb)
        lib.UniLog_RegisterNotifCallback(self._handle, self._notif_callback_ref)


    ##########################################################################
    # Internal sync Logging Method
    
    def _log(self, level: int, msg: str):
        # 1. Capture user caller information 
        caller_info = self._get_caller_info(3)
        
        # 2. Immediate synchronous execution
        self._dispatch_log_to_cgo(level, msg, *caller_info)


    ##########################################################################
    # Internal async Logging Method
    
    async def _async_log(self, level: int, msg: str):
        # 1. Capture user caller information BEFORE moving to background thread
        caller_info = self._get_caller_info(3)
        
        # 2. Offload the blocking C-call to the default thread pool executor (non-blocking)
        await asyncioGetRunningLoop().run_in_executor(
            None, 
            self._dispatch_log_to_cgo, 
            int(level), 
            str(msg), 
            *caller_info
        )


    ##########################################################################
    # Common Logging Core

    # Capture caller metadata from the current stack trace
    def _get_caller_info(self, depth: int):
        caller = inspectStack()[depth]
        filename = osPathBasename(caller.filename)
        lineno = str(caller.lineno)
        function = caller.function
        module = caller.frame.f_globals.get('__name__', 'unknown')
        return filename, lineno, function, module

    # Primary bridge to the Go shared library for all logging events
    def _dispatch_log_to_cgo(self, level, msg, filename, lineno, function, module):
        lib.UniLog_LogWithMetadata(
            self._handle, 
            int(level), 
            str(msg).encode('utf-8'), 
            filename.encode('utf-8'), 
            lineno.encode('utf-8'), 
            function.encode('utf-8'), 
            module.encode('utf-8')
        )

    ##########################################################################
    # Lifecycle and Context management

    # Release the logger session and free shared memory resources
    def close(self):
        """Manually release the logger session and associated resources."""
        if hasattr(self, '_handle') and self._handle:
            lib.UniLog_Close(self._handle)
            self._handle = None


    # Ensure resources are released if the object is garbage collected
    def __del__(self):
        self.close()
    
    ##########################################################################
    # Sync Context Management
    
    # Support for standard 'with' context manager
    def __enter__(self):
        return self

    # Automatic cleanup when exiting 'with' block
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()

    ##########################################################################
    # Async Context Management
    
    # Support for 'async with' context manager
    async def __aenter__(self):
        return self.__enter__()

    # Automatic cleanup when exiting 'async with' block
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        return self.__exit__(exc_type, exc_val, exc_tb)
