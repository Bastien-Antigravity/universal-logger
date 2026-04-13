# Universal Logger: C++ Library

Universal Logger provides a modern C++ RAII wrapper (`UniversalLogger.hpp`) around the Go-based shared library. It is designed to be lightweight, thread-safe, and easy to integrate into existing C++11 projects.

## 🚀 Native Core (DLL)

Since Modernization 2026, the C++ wrapper relies on the **compiled Go shared library** (`libunilog.dll` on Windows).

**Requirement**: Ensure `libunilog.dll` is in your system `PATH` or the same directory as your application.

## 🚀 Features

- **RAII Managed**: Automatic handle cleanup through the `UniLog` class destructor.
- **Header-Only Wrapper**: No separate C++ library to build—just include the header and link the Go library.
- **Telemetry Parity**: Shared library ensures your C++ logs match the format and performance of Go/Rust/Python services.

## 🔧 Installation and Linking

### 1. Build the Go Core
See the root README for Docker-based build instructions to generate `libunilog.dll` and `libunilog.a`.

### 2. Include and Link
```cpp
#include "UniversalLogger.hpp"

int main() {
    // Initialize (loads DLL automatically)
    UniLog logger("standalone", "cpp-app");
    logger.info("C++ is online and powered by Go!");
    return 0;
}
```

#### Compilation (Example with g++)
```bash
g++ -std=c++11 main.cpp -o app -I../libunilog -L../libunilog -lunilog
```

## 📖 Quick Start

### Standard Logging
```cpp
#include "UniversalLogger.hpp"

int main() {
    try {
        UniLog logger("standalone", "demo-app", "standard", UniLog::INFO);

        logger.info("Starting application...");
        logger.debug("Debugging values...");
    } catch (const std::exception& e) {
        std::cerr << "Initialization failed: " << e.what() << std::endl;
    }
    return 0; // RAII cleans up the Go session
}
```

### Automatic Metadata (Macros)
```cpp
#include "UniversalLogger.hpp"

// Captures __FILE__, __LINE__, and __FUNCTION__ automatically
UNILOG_INFO(logger, "Structured logging with full metadata!");
```

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
