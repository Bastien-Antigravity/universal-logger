# Architecture: C++ Library

This document explains the internal design of the C++ facade for the Universal Logger.

## Handle-Based Encapsulation

The C++ library follows the **RAII (Resource Acquisition Is Initialization)** pattern. The `UniversalLogger` class encapsulates the Go-side session handle and ensures that the session is closed when the object goes out of scope.

### 1. Resource Lifecycle
- **Constructor**: Calls `UniLog_Init` from the Go-shared library and stores the returned `uintptr` handle.
- **Destructor**: Automatically calls `UniLog_Close(handle)` to release Go-side resources and prevent memory leaks in the shared core.

```cpp
class UniversalLogger {
    uintptr_t handle;
    // ...
    ~UniversalLogger() {
        if (handle) UniLog_Close(handle);
    }
}
```

### 2. Header-Only Design
To minimize integration friction, the C++ facade is implemented as a thin header-only wrapper (`UniversalLogger.hpp`). This allows developers to include the functionality without needing to build a separate C++ object file/library.

## Callbacks and Function Objects

For configuration updates, the C++ facade supports `std::function<void(const std::string&)>`. This allows for modern C++ lambdas, member functions, or standard function pointers.

### Callback Workflow
1. **Registration**: The user provides a lambda.
2. **Handle Mapping**: The C++ facade registers a static C-style trampoline function with the CGO bridge.
3. **Dispatch**: When the CGO bridge executes the callback, the trampoline identifies the correct `UniversalLogger` instance and invokes the registered `std::function`.

## Concurrency and Thread Safety

The Go core is thread-safe, and the C++ facade maintains this safety. All logging methods (`info`, `debug`, etc.) are safe to call from multiple C++ threads.

### Note on Background Threads
Callbacks for configuration updates are triggered from a **background thread** managed by the Go runtime. C++ developers should use standard synchronization primitives (e.g., `std::mutex`) if they need to update shared C++ state inside the callback lambda.
