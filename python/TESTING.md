# Testing: Python Facade

This document explains how to run the Python test suite for Universal Logger.

## Requirements

Before running the tests, you must build the Go core into a shared library.

```bash
# Build the Go core (creates libunilog.so or .dylib)
make core
```

## Running Tests

We provide three primary test files in the `python/` directory.

### 1. Unified Callback Testing
Verifies the real-time configuration update mechanism (sync and async).
```bash
export PYTHONPATH=$(PWD)/python
export DYLD_LIBRARY_PATH=$(PWD)/libunilog
python3 python/test_unified_callback.py
```

### 2. Async Logging Testing
Verifies non-blocking logging under heavy load.
```bash
python3 python/test_async_logging.py
```

### 3. Basic Functionality Testing
Verifies initialization, log levels, and standard logging.
```bash
python3 python/test_unilog.py
```

## Troubleshooting Tests

### Error: `libunilog shared library not found`
Ensure that the `DYLD_LIBRARY_PATH` (macOS) or `LD_LIBRARY_PATH` (Linux) environment variable is set to the absolute path of the `libunilog/` directory.

### Error: `ModuleNotFoundError: No module named 'unilog'`
Ensure that `PYTHONPATH` includes the `python/` directory so the package can be located.

## Infrastructure Notes

Tests use the `standalone` profile by default. This ensures that no remote configuration server is required to run the automated suite.
