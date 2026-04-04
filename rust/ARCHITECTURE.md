# Architecture: Rust Library

This document explains the internal design of the Rust facade for the Universal Logger.

## Safety and Type-Safety

Traditional FFI (Foreign Function Interface) coding in Rust is inherently `unsafe`. To provide a safe and ergonomic developer experience, this crate uses several techniques to wrap the raw CGO bridge.

## Handling the FFI Boundary

### 1. Pointer Wrapping
All `UniLog` instances are represented as a `uintptr` on the Go side. In Rust, we wrap this in a struct that implements `Drop`.

```rust
pub struct UniLog {
    handle: usize, // Corresponds to Go's uintptr
}

impl Drop for UniLog {
    fn drop(&mut self) {
        unsafe { UniLog_Close(self.handle); }
    }
}
```

This ensures that Go-side memory is automatically cleaned up when the Rust logger goes out of scope.

### 2. Thread Safety
The Rust facade implements `Send` and `Sync` for the `UniLog` struct. This allows the logger to be safely shared across threads, matching the underlying Go core's thread-safe implementation.

## Async and Callbacks

For configuration updates, Rust uses a bridge that translates the C function pointer into a Rust closure.

### Callback Workflow
1. **Registration**: The user provides a closure `Fn(String)`.
2. **Trampoline**: A static "trampoline" function is registered with the CGO bridge.
3. **Dispatch**: When the CGO bridge calls the trampoline, it looks up the original Rust closure and executes it with the serialized JSON update.

## Linking Architecture

The crate uses `rustc-link-lib` to link against `libunilog`. During development, `RUSTFLAGS="-L ..."` is required to inform the linker where the Go-shared library is located. In production, the library is expected to be in a standard system path or relative to the executable.

## Performance Considerations

By utilizing Go's internal non-blocking logging, the Rust facade performs minimal work on the calling thread. The FFI overhead is negligible compared to the network or disk I/O performed by the Go core.
