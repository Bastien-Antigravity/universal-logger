# Universal Logger: Python Facade

High-performance, async-capable Python facade for the Universal Logger. This library provides a native Python experience while leveraging a powerful Go-based core via CGO.

## 🚀 Features

- **Asynchronous First**: Full `asyncio` support for non-blocking logging.
- **Dynamic Config**: `async for` support for real-time configuration updates.
- **Thread Safe**: Safe for multi-threaded and multi-coroutine environments.
- **Easy Deployment**: Standard `setup.py` for pip installation.

## 🔧 Installation

```bash
# Clone the repository
git clone https://github.com/Bastien-Antigravity/universal-logger
cd universal-logger

# Build the Go core first
make core

# Install the Python package
cd python
pip install .
```

## 📖 Quick Start

### Basic Logging
```python
from unilog import UniLog

# Initialize
logger = UniLog(config_profile="standalone", app_name="demo-app")

# Log messages
logger.info("Application started")
logger.debug("Debug information")

# Clean up
logger.close()
```

### Async Logging
```python
import asyncio
from unilog import UniLog

async def main():
    async with UniLog() as logger:
        await logger.async_info("Async logging is easy!")

asyncio.run(main())
```

### Configuration Callbacks (Async)
```python
async with UniLog() as logger:
    async for update in logger.on_config_update():
        print(f"Config changed: {update}")
```

## 🛠️ Configuration Parameters

| Parameter | Default | Description |
| :--- | :--- | :--- |
| `config_profile` | `"standalone"` | `production`, `preprod`, `test`, `standalone` |
| `app_name` | `"python-app"` | Application identifier |
| `logger_profile` | `"standard"` | `standard`, `devel`, `high_perf`, `minimal` |
| `log_level` | `"info"` | `debug`, `info`, `warning`, `error`, `critical` |

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
