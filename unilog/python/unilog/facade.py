#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os.path import basename as osPathBasename
from ctypes import c_int as ctypeC_int
from sys import _getframe as sysGetFrame
from json import loads as jsonLoads
from asyncio import get_running_loop as asyncioGetRunningLoop
from typing import Union, Dict, List, Callable, Optional, Set, Any, Tuple

from .models import LogLevel
from .listeners import ConfigUpdateListener
from .lib_loader import lib, CALLBACK_TYPE


class UniLog:
    """
    Python Facade for the Universal Logger (Go) shared library.
    Provides integrated configuration management and high-performance logging.
    """

    def __init__(self, app_name: str = "python-app", config_profile: str = "standalone", 
                 logger_profile: str = "standard", log_level: Union[str, int, LogLevel] = "info", 
                 use_local_notifier: bool = False, config_handle: int = 0) -> None:
        if not lib:
            raise RuntimeError("libunilog shared library not found. Please ensure it is built.")
        
        # Convert string log level to int for Go
        level_val = LogLevel.from_str(log_level) if isinstance(log_level, str) else int(log_level)

        self._handle: int = lib.UniLog_Init(
            app_name.encode('utf-8'), 
            config_profile.encode('utf-8'), 
            logger_profile.encode('utf-8'), 
            ctypeC_int(level_val),
            ctypeC_int(1 if use_local_notifier else 0),
            config_handle
        )
        self._callback_ref: Optional[Any] = None # Keep reference to avoid GC
        self._sync_subscribers: List[Callable[[Dict[str, Any]], None]] = []
        self._async_listeners: Set[ConfigUpdateListener] = set()
        self._initialized_bridge: bool = False


    ##########################################################################
    # Logging Methods
    
    def debug(self, msg: str) -> None: self._log(LogLevel.DEBUG, msg)
    def info(self, msg: str) -> None: self._log(LogLevel.INFO, msg)
    def warning(self, msg: str) -> None: self._log(LogLevel.WARNING, msg)
    def error(self, msg: str) -> None: self._log(LogLevel.ERROR, msg)
    def critical(self, msg: str) -> None: self._log(LogLevel.CRITICAL, msg)

    async def async_debug(self, msg: str) -> None: await self._async_log(LogLevel.DEBUG, msg)
    async def async_info(self, msg: str) -> None: await self._async_log(LogLevel.INFO, msg)
    async def async_warning(self, msg: str) -> None: await self._async_log(LogLevel.WARNING, msg)
    async def async_error(self, msg: str) -> None: await self._async_log(LogLevel.ERROR, msg)
    async def async_critical(self, msg: str) -> None: await self._async_log(LogLevel.CRITICAL, msg)

    # Specialized Domain Methods
    def logon(self, msg: str) -> None: self._log(LogLevel.LOGON, msg)
    def logout(self, msg: str) -> None: self._log(LogLevel.LOGOUT, msg)
    def trade(self, msg: str) -> None: self._log(LogLevel.TRADE, msg)
    def schedule(self, msg: str) -> None: self._log(LogLevel.SCHEDULE, msg)
    def report(self, msg: str) -> None: self._log(LogLevel.REPORT, msg)
    def stream(self, msg: str) -> None: self._log(LogLevel.STREAM, msg)

    # Log level accessors
    def set_level(self, level: Union[str, int, LogLevel]) -> None:
        """Change the current log level dynamically."""
        if isinstance(level, str):
            level = LogLevel.from_str(level)
        lib.UniLog_SetLevel(self._handle, int(level))
        
    def get_level(self) -> LogLevel:
        """Retrieve the current log level from the Go core."""
        return LogLevel(lib.UniLog_GetLevel(self._handle))
        

    ##########################################################################
    # Metadata Methods ---

    def add_metadata(self, key: str, value: str) -> None:
        """Add a single key-value pair to all future logs."""
        lib.UniLog_AddMetadata(self._handle, key.encode('utf-8'), value.encode('utf-8'))

    def set_metadata(self, metadata: Dict[str, Any]) -> None:
        """Replace all existing metadata with the provided dictionary."""
        import json
        json_data = json.dumps(metadata)
        lib.UniLog_SetMetadata(self._handle, json_data.encode('utf-8'))

    ##########################################################################
    # Config Methods ---

    def get_config(self, section: str, key: str, default: Optional[str] = None) -> Optional[str]:
        """Retrieve a configuration value from the distributed config service (Zero-Leak FFI)."""
        res_ptr = lib.UniLog_Config_Get(self._handle, section.encode('utf-8'), key.encode('utf-8'))
        if not res_ptr:
            return default
        try:
            from ctypes import string_at
            return string_at(res_ptr).decode('utf-8')
        finally:
            if hasattr(lib, "DistConf_FreeString"):
                lib.DistConf_FreeString(res_ptr)

    def set_config(self, section: str, key: str, value: str) -> None:
        """Update a configuration value in the memory configuration."""
        lib.UniLog_Config_Set(self._handle, section.encode('utf-8'), key.encode('utf-8'), value.encode('utf-8'))

    # Trigger on_config_update regarding the caller and caller method
    def _dispatch_update(self, json_data: bytes) -> None:
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

    def on_config_update(self, callback: Optional[Callable[[Dict[str, Any]], None]] = None) -> Optional[ConfigUpdateListener]:
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
            lib.UniLog_OnConfigUpdate(self._handle, self._callback_ref)
            print("!!! Python: UniLog_OnConfigUpdate call finished.")
            self._initialized_bridge = True

        if callback is not None:
            self._sync_subscribers.append(callback)
            return None
        
        return ConfigUpdateListener(self)


    ##########################################################################
    # Local Notifier Methods

    def on_notification(self, callback: Callable[[Dict[str, Any]], None]) -> None:
        """
        Registers a callback for local notifications.
        The callback will receive a dictionary parsed from the notification JSON.
        """
        def _bridge_cb(json_data: bytes) -> None:
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
    
    def _log(self, level: Union[int, LogLevel], msg: str) -> None:
        # 1. Capture user caller information 
        caller_info = self._get_caller_info(3)
        
        # 2. Immediate synchronous execution
        self._dispatch_log_to_cgo(level, msg, *caller_info)


    ##########################################################################
    # Internal async Logging Method
    
    async def _async_log(self, level: Union[int, LogLevel], msg: str) -> None:
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
    def _get_caller_info(self, depth: int) -> Tuple[str, str, str, str]:
        frame = sysGetFrame(depth)
        filename = osPathBasename(frame.f_code.co_filename)
        lineno = str(frame.f_lineno)
        function = frame.f_code.co_name
        module = frame.f_globals.get('__name__', 'unknown')
        return filename, lineno, function, module

    # Primary bridge to the Go shared library for all logging events
    def _dispatch_log_to_cgo(self, level: Union[int, LogLevel], msg: str, filename: str, lineno: str, function: str, module: str) -> None:
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
    def close(self) -> None:
        """Manually release the logger session and associated resources."""
        if hasattr(self, '_handle') and self._handle:
            lib.UniLog_Close(self._handle)
            self._handle = None


    # Ensure resources are released if the object is garbage collected
    def __del__(self) -> None:
        self.close()
    
    ##########################################################################
    # Sync Context Management
    
    # Support for standard 'with' context manager
    def __enter__(self) -> 'UniLog':
        return self

    # Automatic cleanup when exiting 'with' block
    def __exit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        self.close()

    ##########################################################################
    # Async Context Management
    
    # Support for 'async with' context manager
    async def __aenter__(self) -> 'UniLog':
        return self.__enter__()

    # Automatic cleanup when exiting 'async with' block
    async def __aexit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        return self.__exit__(exc_type, exc_val, exc_tb)
