# Architecture: Python Facade

This document explains the internal design of the Python library for Universal Logger.

## Design Philosophy

The Python library is designed to be **thin, idiomatic, and non-blocking**. It uses `ctypes` to communicate with the shared Go core and `asyncio` to provide a modern developer experience.

## Package Structure

The core logic is contained within the `unilog/` internal package:
- **`facade.py`**: The primary `UniLog` class. Handles initialization and dispatches calls to the bridge.
- **`lib_loader.py`**: Locates and loads the `libunilog` shared library (`.so`, `.dylib`, or `.dll`).
- **`listeners.py`**: Implements the `ConfigUpdateListener` for async iteration.
- **`models.py`**: Defines the `LogLevel` enumeration for type checking.

## Configuration Callback Mechanism

Universal Logger supports real-time configuration updates. Because these updates originate from a background Go goroutine, we bridge them into the Python event loop using a thread-safe pattern.

### 1. Callback Registration
When `on_config_update()` is called, Python registers a `CFUNCTYPE` callback with the CGO bridge.

### 2. Async Dispatching
When a configuration update is received:
1.  **Background Thread**: The CGO bridge calls the registered C function in a background thread.
2.  **Queue Injection**: The `UniLog._dispatch_update` method (Python) acquires the GIL and pushes the new configuration into a thread-safe `asyncio.Queue`.
3.  **Event Loop**: The `ConfigUpdateListener` (running in the main event loop) yields the data from the queue to the `async for` loop.

## Async Logging Architecture

Logging is notoriously blocking if I/O is involved. To keep the Python event loop responsive, `async_info`, `async_debug`, etc., use a thread pool executor:

```python
# From facade.py
async def _async_log(self, level, msg):
    caller_info = self._get_caller_info(3) # Capture stack BEFORE moving to background
    await asyncio.get_running_loop().run_in_executor(
        None, 
        self._dispatch_log_to_cgo, 
        level, msg, *caller_info
    )
```

This ensures that even if the Go-side logging engine is momentarily under heavy load, the Python application remains snappy.

## Resource Lifecycle

The `UniLog` class implements both the `with` and `async with` context manager protocols:
- **`__enter__` / `__aenter__`**: Returns the logger instance.
- **`__exit__` / `__aexit__`**: Calls `self.close()` to release the Go-side session handle.
- **`__del__`**: A safety fallback to release handles if the object is garbage collected without being closed.
