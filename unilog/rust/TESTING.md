# Testing: Rust Library

This document explains how to run the Rust test suite for Universal Logger.

## Requirements

Before running the tests, you must build the Go core into a shared library.

```bash
# Build the Go core (creates libunilog.so or .dylib)
make core
```

## Running Tests

We use `cargo test` to verify the Rust facade's functionality.

### 1. Verification of Safe Pointers
Tests that the facade correctly initializes and cleans up Go handles.
```bash
export RUSTFLAGS="-L $(pwd)/libunilog"
cargo test
```

### 2. Runtime Library Discovery
Ensure that the dynamic linker can find the Go shared library at runtime.

```bash
# macOS
export DYLD_LIBRARY_PATH=$(pwd)/libunilog
cargo test

# Linux
export LD_LIBRARY_PATH=$(pwd)/libunilog
cargo test
```

## Troubleshooting Tests

### Error: `linker 'cc' failed: exit code: 1`
This usually means `cargo` cannot find the `libunilog` shared library at link time. Ensure that `RUSTFLAGS="-L ..."` correctly points to the directory containing `libunilog.so` or `libunilog.dylib`.

### Error: `Library not loaded: @rpath/libunilog.dylib`
This is a runtime error. Ensure that `DYLD_LIBRARY_PATH` (macOS) or `LD_LIBRARY_PATH` (Linux) is set correctly.

## Infrastructure Notes

Rust tests are integrated into the root `Makefile` via `make rust`. This is the recommended way to run the full cross-platform verification.
