#!/usr/bin/env python
# coding:utf-8
"""
=============================================================================
ESSENTIAL PROCESS:
Primary Python facade for the universal-logger library, providing synchronous
and asynchronous logging methods, caller stack inspection, and CGO bridge interop.

DATA FLOW:
1. Input: Application log messages, log levels, metadata, and config keys.
2. Logic: Inspects caller frames via sys._getframe(), serializes metadata to JSON,
   and delegates calls to the libunilog CGO shared library.
3. Output: Dispatches logs, config updates, and local alerts across FFI boundaries.

KEY PARAMETERS:
- app_name: Identifier of the client microservice or application.
- config_profile: Configuration strategy ('standalone' or 'service').
- logger_profile: Logging sink configuration profile ('devel', 'production').
=============================================================================
"""

from asyncio import get_running_loop as asyncioGetRunningLoop
from ctypes import c_int as ctypeC_int, string_at as ctypeStringAt
from json import dumps as jsonDumps, loads as jsonLoads
from os.path import basename as osPathBasename
from sys import _getframe as sysGetFrame, stderr as sysStderr
from typing import Any as typeAny, Callable as typeCallable, Dict as typeDict, List as typeList, Optional as typeOptional, Set as typeSet, Tuple as typeTuple, Union as typeUnion

from .models import LogLevel
from .listeners import ConfigUpdateListener
from .lib_loader import lib, CALLBACK_TYPE

# -----------------------------------------------------------------------------------------------


class UniLog:
    """
    Python Facade for the Universal Logger (Go) shared library.
    Provides integrated configuration management and high-performance logging.
    """

    Name = "UniLog"

    # -----------------------------------------------------------------------------------------------

    def __init__(self, app_name: str = "python-app", config_profile: str = "standalone",
                 logger_profile: str = "standard", log_level: typeUnion[str, int, LogLevel] = "info",
                 use_local_notifier: bool = False, config_handle: int = 0) -> None:
        if not lib:
            raise RuntimeError("libunilog shared library not found. Please ensure it is built.")

        # Convert string log level to int for Go
        level_val = LogLevel.from_str(log_level) if isinstance(log_level, str) else int(log_level)

        self._handle: typeOptional[int] = lib.UniLog_Init(
            app_name.encode('utf-8'),
            config_profile.encode('utf-8'),
            logger_profile.encode('utf-8'),
            ctypeC_int(int(level_val)),
            ctypeC_int(1 if use_local_notifier else 0),
            config_handle
        )
        self._min_level: int = int(level_val)
        self._callback_ref: typeOptional[typeAny] = None  # Keep reference to avoid GC
        self._notif_callback_ref: typeOptional[typeAny] = None
        self._sync_subscribers: typeList[typeCallable[[typeDict[str, typeAny]], None]] = []
        self._async_listeners: typeSet[ConfigUpdateListener] = set()
        self._initialized_bridge: bool = False

    # -----------------------------------------------------------------------------------------------

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

    # -----------------------------------------------------------------------------------------------

    def set_level(self, level: typeUnion[str, int, LogLevel]) -> None:
        """Change the current log level dynamically."""
        if isinstance(level, str):
            level = LogLevel.from_str(level)
        lvl_int = int(level)
        self._min_level = lvl_int
        lib.UniLog_SetLevel(self._handle, lvl_int)

    # -----------------------------------------------------------------------------------------------

    def get_level(self) -> LogLevel:
        """Retrieve the current log level from the Go core."""
        return LogLevel(lib.UniLog_GetLevel(self._handle))

    # -----------------------------------------------------------------------------------------------

    def add_metadata(self, key: str, value: str) -> None:
        """Add a single key-value pair to all future logs."""
        lib.UniLog_AddMetadata(self._handle, key.encode('utf-8'), value.encode('utf-8'))

    # -----------------------------------------------------------------------------------------------

    def set_metadata(self, metadata: typeDict[str, typeAny]) -> None:
        """Replace all existing metadata with the provided dictionary."""
        json_data = jsonDumps(metadata)
        lib.UniLog_SetMetadata(self._handle, json_data.encode('utf-8'))

    # -----------------------------------------------------------------------------------------------

    def get_config(self, section: str, key: str, default: typeOptional[str] = None) -> typeOptional[str]:
        """Retrieve a configuration value from the distributed config service (Zero-Leak FFI)."""
        res_ptr = lib.UniLog_Config_Get(self._handle, section.encode('utf-8'), key.encode('utf-8'))
        if not res_ptr:
            return default
        try:
            return ctypeStringAt(res_ptr).decode('utf-8')
        finally:
            if hasattr(lib, "DistConf_FreeString"):
                lib.DistConf_FreeString(res_ptr)

    # -----------------------------------------------------------------------------------------------

    def set_config(self, section: str, key: str, value: str) -> None:
        """Update a configuration value in the memory configuration."""
        lib.UniLog_Config_Set(self._handle, section.encode('utf-8'), key.encode('utf-8'), value.encode('utf-8'))

    # -----------------------------------------------------------------------------------------------

    def _dispatch_update(self, json_data: bytes) -> None:
        """Internal bridge called from Go shared library background thread."""
        try:
            raw_val = json_data.decode('utf-8')
            data = jsonLoads(raw_val)

            logger_section = data.get("logger") if isinstance(data, dict) else None
            if isinstance(logger_section, dict) and "level" in logger_section:
                self._min_level = int(LogLevel.from_str(str(logger_section["level"])))

            # 1. Dispatch to synchronous subscribers
            for cb in self._sync_subscribers:
                try:
                    cb(data)
                except Exception as err:
                    sysStderr.write(f"UniLog : subscriber callback failed ({err})\n")

            # 2. Dispatch to asynchronous listeners (thread-safe)
            for listener in list(self._async_listeners):  # Copy list to avoid concurrent mutation
                listener._put(data)
        except Exception as err:
            sysStderr.write(f"UniLog : failed to dispatch config update ({err})\n")

    # -----------------------------------------------------------------------------------------------

    def on_config_update(self, callback: typeOptional[typeCallable[[typeDict[str, typeAny]], None]] = None) -> typeOptional[ConfigUpdateListener]:
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
            lib.UniLog_OnConfigUpdate(self._handle, self._callback_ref)
            self._initialized_bridge = True

        if callback is not None:
            self._sync_subscribers.append(callback)
            return None

        return ConfigUpdateListener(self)

    # -----------------------------------------------------------------------------------------------

    def on_notification(self, callback: typeCallable[[typeDict[str, typeAny]], None]) -> None:
        """
        Registers a callback for local notifications.
        The callback will receive a dictionary parsed from the notification JSON.
        """
        def _bridge_cb(json_data: bytes) -> None:
            try:
                data = jsonLoads(json_data.decode('utf-8'))
                callback(data)
            except Exception as err:
                sysStderr.write(f"UniLog : notification callback failed ({err})\n")

        # Keep a reference to the bridge callback to avoid GC
        self._notif_callback_ref = CALLBACK_TYPE(_bridge_cb)
        lib.UniLog_RegisterNotifCallback(self._handle, self._notif_callback_ref)

    # -----------------------------------------------------------------------------------------------

    def _log(self, level: typeUnion[int, LogLevel], msg: str) -> None:
        lvl_int = int(level)
        if lvl_int < self._min_level:
            return
        # 1. Capture user caller information
        caller_info = self._get_caller_info(3)

        # 2. Immediate synchronous execution
        self._dispatch_log_to_cgo(lvl_int, msg, *caller_info)

    # -----------------------------------------------------------------------------------------------

    async def _async_log(self, level: typeUnion[int, LogLevel], msg: str) -> None:
        lvl_int = int(level)
        if lvl_int < self._min_level:
            return
        # 1. Capture user caller information BEFORE moving to background thread
        caller_info = self._get_caller_info(3)

        # 2. Offload the blocking C-call to the default thread pool executor (non-blocking)
        await asyncioGetRunningLoop().run_in_executor(
            None,
            self._dispatch_log_to_cgo,
            lvl_int,
            str(msg),
            *caller_info
        )

    # -----------------------------------------------------------------------------------------------

    # Capture caller metadata from the current stack trace
    def _get_caller_info(self, depth: int) -> typeTuple[str, str, str, str]:
        try:
            frame = sysGetFrame(depth)
            filename = osPathBasename(frame.f_code.co_filename)
            lineno = str(frame.f_lineno)
            function = frame.f_code.co_name
            module = frame.f_globals.get('__name__', 'unknown')
            return filename, lineno, function, module
        except (ValueError, AttributeError):
            return "unknown", "0", "unknown", "unknown"

    # -----------------------------------------------------------------------------------------------

    # Primary bridge to the Go shared library for all logging events
    def _dispatch_log_to_cgo(self, level: typeUnion[int, LogLevel], msg: str, filename: str, lineno: str, function: str, module: str) -> None:
        lib.UniLog_LogWithMetadata(
            self._handle,
            int(level),
            str(msg).encode('utf-8'),
            filename.encode('utf-8'),
            lineno.encode('utf-8'),
            function.encode('utf-8'),
            module.encode('utf-8')
        )

    # -----------------------------------------------------------------------------------------------

    # Release the logger session and free shared memory resources
    def close(self) -> None:
        """Manually release the logger session and associated resources."""
        if hasattr(self, '_handle') and self._handle:
            lib.UniLog_Close(self._handle)
            self._handle = None

    # -----------------------------------------------------------------------------------------------

    def __del__(self) -> None:
        self.close()

    # -----------------------------------------------------------------------------------------------

    def __enter__(self) -> 'UniLog':
        return self

    # -----------------------------------------------------------------------------------------------

    def __exit__(self, exc_type: typeAny, exc_val: typeAny, exc_tb: typeAny) -> None:
        self.close()

    # -----------------------------------------------------------------------------------------------

    async def __aenter__(self) -> 'UniLog':
        return self.__enter__()

    # -----------------------------------------------------------------------------------------------

    async def __aexit__(self, exc_type: typeAny, exc_val: typeAny, exc_tb: typeAny) -> None:
        return self.__exit__(exc_type, exc_val, exc_tb)
