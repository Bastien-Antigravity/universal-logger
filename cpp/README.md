# Universal Logger: C++ Library

Universal Logger provides a modern C++ wrapper (`UniversalLogger.hpp`) around the Go-based shared library. It is designed to be lightweight, thread-safe, and easy to integrate into existing C++11 projects.

## 🚀 Features

- **RAII Managed**: Automatic handle cleanup through the `UniversalLogger` class destructor.
- **Header-Only Wrapper**: No separate C++ library to build—just include the header and link the Go library.
- **Cross-Platform**: Compatible with standard C++ compilers (g++, clang++, msvc).

## 🔧 Installation and Linking

### 1. Build the Go Core
```bash
make core
```

### 2. Include and Link
The C++ library depends on the `libunilog` shared library and its generated C header.

```cpp
#include "UniversalLogger.hpp"

int main() {
    UniversalLogger logger("standalone", "cpp-app");
    logger.info("C++ is online!");
    return 0;
}
```

#### Compilation (Example with g++)
```bash
g++ -std=c++11 main.cpp -o app -I./libunilog -L./libunilog -lunilog
```

## 📖 Quick Start

```cpp
#include "UniversalLogger.hpp"

int main() {
    // Basic Initialization
    UniversalLogger logger("standalone", "demo-app", "devel", LogLevel::INFO);

    // Standard Logging
    logger.info("Starting application...");
    logger.debug("Debugging values...");
    
    // Dynamic Config Updates
    logger.on_config_update([](const std::string& update) {
        std::cout << "Config update received: " << update << std::endl;
    });

    return 0; // RAII cleans up the Go session
}
```

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
