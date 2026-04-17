# Universal Logger: Python Facade

High-performance, async-capable Python facade for the Universal Logger. This library provides a native Python experience while leveraging a powerful Go-based core via CGO.

## 🚀 Native Core (DLL)

Since Modernization 2026, the Python facade relies on the **compiled Go shared library** (`libunilog.dll` on Windows).

**Requirement**: Ensure `libunilog.dll` is in your system `PATH` or the same directory as your application.

## 🚀 Features

- **Asynchronous First**: Full `asyncio` support for non-blocking logging.
- **Dynamic Config**: `async for` support for real-time configuration updates.
- **Thread Safe**: Safe for multi-threaded and multi-coroutine environments.
- **Centralized Engine**: Shares the exact same logic and performance as the Rust and Go clients.

## 🔧 Installation

```bash
# 1. Ensure the core is built (generates libunilog.dll)
# See root README for Docker-based build instructions

# 2. Install the Python package
cd universal-logger/python
pip install .
```

## 📖 Quick Start

### Basic Logging
```python
from unilog import UniLog

# Initialize (loads libunilog.dll automatically)
logger = UniLog(app_name="demo-app", config_profile="standalone")

# Log messages
logger.info("Application started")
logger.debug("Debug information")

# Clean up
logger.close()
```

### Async Logging (Recommended)
```python
import asyncio
from unilog import UniLog

async def main():
    async with UniLog() as logger:
        await logger.async_info("Async logging is powered by Go!")

asyncio.run(main())
```

## 🛠️ Configuration Parameters

| Parameter | Default | Description |
| :--- | :--- | :--- |
| `app_name` | `"python-app"` | Application identifier (Trimmmed for exact match) |
| `config_profile` | `"standalone"` | `production`, `preprod`, `test`, `standalone` |
| `logger_profile` | `"standard"` | `standard`, `devel`, `high_perf`, `minimal` |
| `log_level` | `"info"` | `debug`, `info`, `warning`, `error`, `critical` |
| `use_local_notifier` | `False` | Enable in-process 1024-buffered notification queue |

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
