# Universal Logger: Rust Library

Safe, high-performance, and type-safe Rust idiomatic facade for the Universal Logger. This crate provides a native Rust experience while leveraging a powerful Go-based core via CGO.

## 🚀 Native Core (DLL)

Since Modernization 2026, the Rust facade relies on the **compiled Go shared library** (`libunilog.dll` on Windows).

**Requirement**: Ensure `libunilog.dll` is in your system `PATH` or the same directory as your application.

## 🚀 Features

- **Safe Pointers**: Handle memory securely with Rust lifetimes and `Drop` trait for automated cleanup.
- **Cargo Integration**: Full support for your Rust builds with automated linking via `build.rs`.
- **Async Ready**: Seamless integration with `tokio` or other async runtimes via callbacks.
- **Zero Overhead**: Minimal abstraction over the low-level FFI calls.

## 🔧 Installation

Add Universal Logger to your `Cargo.toml`:

```toml
[dependencies]
unilog-rs = { path = "../rust" }
```

## 📖 Quick Start

### Basic Usage
```rust
use unilog_rs::{UniLog, LogLevel};

fn main() {
    // 1. Initialize the logger session (loads libunilog.dll)
    let logger = UniLog::new(
        "standalone",   // config_profile
        "rust-demo",     // app_name
        "standard",     // logger_profile
        LogLevel::Info, // log_level
        false           // use_local_notifier
    ).expect("Failed to initialize logger");

    // 2. Log messages
    logger.info("Rust is online and powered by Go!");
    logger.debug("Debug message");
}
```

### Automatic Metadata (Macros)
```rust
use unilog_rs::{unilog_info, unilog_error};

// Captures file!(), line!(), and module_path!() automatically
unilog_info!(logger, "Structured logging with metadata!");
```

## 🛠️ Linking Requirements

The `unilog-rs` crate includes a `build.rs` that automatically searches for `libunilog.a` (import library) in the `../libunilog` directory.

1. **Build the Go Core**: See root README for Docker-based build instructions.
2. **Environment**: On Windows, add the directory containing `libunilog.dll` to your `PATH`.

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
