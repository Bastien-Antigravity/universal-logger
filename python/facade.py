import os
import ctypes
import json
import inspect
from pathlib import Path

# Load matching Go-compiled shared library
lib_path = Path(__file__).parent / "libdistconf_flexlog.so"
if not lib_path.exists():
    # Try .dylib for macOS or .dll for Windows during development
    for ext in [".dylib", ".dll"]:
        if (lib_path.with_suffix(ext)).exists():
            lib_path = lib_path.with_suffix(ext)
            break

lib = ctypes.CDLL(str(lib_path))

# Define C-API types
lib.NewFacade.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
lib.NewFacade.restype = ctypes.c_void_p

lib.FreeFacade.argtypes = [ctypes.c_void_p]

lib.LogWithMetadataC.argtypes = [
    ctypes.c_void_p, ctypes.c_int, ctypes.c_char_p, 
    ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p
]

CALLBACK_TYPE = ctypes.CFUNCTYPE(None, ctypes.c_char_p)
lib.RegisterUpdateCallback.argtypes = [ctypes.c_void_p, CALLBACK_TYPE]

class LogLevel:
    NOTSET = 0
    DEBUG = 10
    INFO = 20
    WARNING = 30
    ERROR = 40
    CRITICAL = 50

class DistconfFlexlogFacade:
    def __init__(self, config_profile="standalone", app_name="python-app", 
                 logger_profile="standard", log_level="info"):
        self._handle = lib.NewFacade(
            config_profile.encode(), 
            app_name.encode(), 
            logger_profile.encode(), 
            log_level.encode()
        )
        self._callback_ref = None # Keep reference to avoid GC

    def close(self):
        if self._handle:
            lib.FreeFacade(self._handle)
            self._handle = None

    def __del__(self):
        self.close()

    def _log(self, level, msg):
        # Capture stack info
        caller = inspect.stack()[2]
        filename = os.path.basename(caller.filename)
        lineno = str(caller.lineno)
        function = caller.function
        module = caller.frame.f_globals.get('__name__', 'unknown')

        lib.LogWithMetadataC(
            self._handle, 
            level, 
            msg.encode(), 
            filename.encode(), 
            lineno.encode(), 
            function.encode(), 
            module.encode()
        )

    def debug(self, msg): self._log(LogLevel.DEBUG, msg)
    def info(self, msg): self._log(LogLevel.INFO, msg)
    def warning(self, msg): self._log(LogLevel.WARNING, msg)
    def error(self, msg): self._log(LogLevel.ERROR, msg)
    def critical(self, msg): self._log(LogLevel.CRITICAL, msg)

    def on_config_update(self, callback):
        """
        Register a Python callback for configuration updates.
        The callback will receive a dictionary of the updated configuration.
        """
        def wrapped_callback(json_data):
            try:
                data = json.loads(json_data.decode())
                callback(data)
            except Exception as e:
                print(f"Error in config update callback: {e}")

        self._callback_ref = CALLBACK_TYPE(wrapped_callback)
        lib.RegisterUpdateCallback(self._handle, self._callback_ref)
