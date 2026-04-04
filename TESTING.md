# Testing: Universal Logger

This document provides a high-level overview of the testing strategy for the `universal-logger` project. Since this is a multi-language project, the test suite is divided by platform, with a focus on verifying the CGO bridge integrity.

## Testing Philosophy

Our testing strategy follows the **Contract-First** approach:
1. **Core Logic**: The Go core (`src/`) is tested for orchestration and synergy.
2. **Bridge Integrity**: The CGO bridge is verified by ensuring that handles and callbacks are correctly passed across the FFI boundary.
3. **Facade Idioms**: Each language facade (Python, Rust, C++) is tested to ensure it behaves like a native library while correctly interacting with the shared core.

## Running Tests by Platform

### 1. Unified Makefile (Convenience)
The root `Makefile` provides a quick way to run all primary tests:

```bash
# Run Go core tests
make core_tests

# Run Python-specific tests
make python

# Build and run C++ tests
make cpp

# Build and run Rust tests
make rust
```

### 2. Go Core (`src/`)
Automated tests verify the internal state of the facade and the capability mapping between subsystems.
```bash
go test -v ./src/...
```

### 3. Python (`python/`)
Verifies asynchronous logging, thread-safe configuration updates, and the `async for` listener.
```bash
export PYTHONPATH=$(PWD)/python
python3 python/test_unilog.py
python3 python/test_unified_callback.py
python3 python/test_async_logging.py
```

### 4. Rust (`rust/`)
Verifies the safe pointers and memory management.
```bash
cd rust && cargo test
```

### 5. C++ (`cpp/`)
Verifies the C++ RAII wrapper and basic logging functionality.
```bash
make -C cpp
./cpp/unilog_cpp
```

## Cross-Platform Verification Matrix

Every release is verified against the following matrix:

| Platform | Feature | Test Target |
| :--- | :--- | :--- |
| **Go** | Orchestration | `src/facade/facade_test.go` |
| **Python** | Async I/O | `python/test_async_logging.py` |
| **Python** | Callbacks | `python/test_unified_callback.py` |
| **Rust** | Memory Safety | `rust/src/lib.rs` (Doc tests) |
| **C++** | Thread Safety | `cpp/main.cpp` |
| **VBA** | Message Pump | Manual (Excel Mock) |

## Infrastructure Requirements

- **Local Development**: Use the `standalone` profile to run tests without requiring a remote `distributed-config` server.
- **Environment Variables**: Ensure `DYLD_LIBRARY_PATH` (macOS) or `LD_LIBRARY_PATH` (Linux) points to the `libunilog/` directory during execution.
