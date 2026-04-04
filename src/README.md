# Universal Logger: Go Core and CGO Bridge

The Go core is the heart of the Universal Logger. It orchestrates the configuration and logging subsystems and exposes them via a high-performance C ABI using Go's `c-shared` build mode.

## 🚀 Key Responsibilities

- **System Orchestration**: Coordinates `distributed-config` and `flexible-logger`.
- **CGO Bridge**: Provides a stable, low-latency FFI boundary for other languages.
- **Session Handling**: Manages the lifecycle of multiple independent logger instances.
- **Multithreaded Safety**: Ensures that shared resources are protected via Go's concurrency primitives.

## 🔧 Building the Core

To expose the Go core to other languages, it must be compiled into a dynamic shared library.

```bash
# General build command
go build -buildmode=c-shared -o libunilog/libunilog.so src/cgo_bridge/*.go

# macOS specific (build dylib)
go build -buildmode=c-shared -o libunilog/libunilog.dylib src/cgo_bridge/*.go
```

The build process generates two files in the `libunilog/` directory:
1.  **Shared Library**: `libunilog.so` (Linux) or `libunilog.dylib` (macOS).
2.  **C Header**: `libunilog.h` (used by C/C++, Rust, and Python via `ctypes`).

## 🛠️ Components

- **`src/bootstrap/`**: Handles the alignment of config discovery data and logger sinks.
- **`src/cgo_bridge/`**: Contains the `//export` functions that define the C ABI.
- **`src/config/`**: Thin wrapper around the `distributed-config` library.
- **`src/logger/`**: Implementation of the `UniLog` facade and engine management.

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed Go-side test instructions.
