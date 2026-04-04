# Universal Logger: Rust Library

Safe, high-performance, and type-safe Rust idiomatic facade for the Universal Logger. This crate provides a native Rust experience while leveraging a powerful Go-based core via CGO.

## 🚀 Features

- **Safe Pointers**: Handle memory securely with Rust lifetimes and `Box`.
- **Cargo Integration**: Full support for your Rust builds.
- **Async Ready**: Seamless integration with `tokio` or other async runtimes.
- **Zero Overhead**: Minimal abstraction over the low-level FFI calls.

## 🔧 Installation

Add Universal Logger to your `Cargo.toml`:

```toml
[dependencies]
unilog-rs = { path = "../rust" }
```

## 📖 Quick Start

```rust
use unilog_rs::{UniLog, LogLevel};

fn main() {
    // Initialize the logger
    let logger = UniLog::builder()
        .config_profile("standalone")
        .app_name("rust-demo")
        .build()
        .expect("Failed to initialize logger");

    // Log messages
    logger.info("Rust is online!");
    logger.debug("Debug message");
}
```

## 🛠️ Linking Requirements

Since this crate calls a Go-shared library, you must ensure the library is available at link time and runtime.

### 1. Build the Go Core
```bash
make core
```

### 2. Set Linker Flags
When building your Rust app, tell `cargo` where to find the library:
```bash
export RUSTFLAGS="-L $(pwd)/libunilog"
cargo build
```

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
