# Distconf-Flexlog Facade

A unified high-performance facade that orchestrates [Distributed-Config](https://github.com/Bastien-Antigravity/distributed-config) and [Flexible-Logger](https://github.com/Bastien-Antigravity/flexible-logger).

## 🏗️ Architecture

Following the **SystemFacade** pattern, this library decouples application logic from the underlying configuration and logging drivers.

```text
       +---------------------------------------+
       |           Application Logic           |
       +---------------------------------------+
                     |
                     v
       +---------------------------------------+
       |        DistconfFlexlogFacade          |
       |     (src/facade/facade.go)            |
       +---------------------------------------+
          /                         \
         v                           v
+-------------------+       +-------------------+
| Distributed Config|       |  Flexible Logger  |
| (Sync/Async REST) |       | (Binary/Capnp/UDP)|
+-------------------+       +-------------------+
```

## 🚀 Quick Start

Initialize the entire system with a single call using the `factory`:

```go
import (
    "github.com/Bastien-Antigravity/distconf-flexlog/src/factory"
    "github.com/Bastien-Antigravity/distconf-flexlog/src/models"
)

func main() {
    params := models.MFacadeParams{
        ConfigProfile: "standalone",
        AppName:       "my-service",
        LoggerProfile: "devel",
        LogLevel:      "debug",
    }

    facade := factory.NewFacade(params)
    defer facade.Close()

    facade.Info("System is online!")
}
```

## 📊 Operational Scenarios

| Scenario | Config Profile | Logger Profile | Primary Use Case |
| :--- | :--- | :--- | :--- |
| **1: Local Dev** | `standalone` | `devel` | Standard local development with text logs. |
| **2: Production** | `production` | `standard` | Full remote config + Multi-sink (Network/File) logging. |
| **3: High Load** | `production` | `high_perf` | Optimized for high-throughput async network logging. |
| **4: Testing** | `test` | `minimal` | Lightweight setup for CI/CD and automated tests. |
| **5: Monitor** | `preprod` | `notif_logger` | Notifications focus with remote config monitoring. |
| **6: Critical** | `production` | `no_lock` | Ultra-low latency using lock-free atomic logging. |

## 🛠️ Configuration Parameters

Initialize using the `MFacadeParams` struct:

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `ConfigProfile` | `string` | `production`, `preprod`, `test`, `standalone` |
| `AppName` | `string` | Application identifier (used for both systems) |
| `LoggerProfile` | `string` | `standard`, `devel`, `high_perf`, `minimal`, `notif_logger`, `no_lock` |
| `LogLevel` | `string` | `debug`, `info`, `warning`, `error`, `critical` |
| `PublicIP` | `string` | Optional: used for remote identification (defaults to 127.0.0.1) |

## 🛠️ Internal Maintenance

This facade includes alignment fixes for:
*   **Field Mapping**: Links `distributed-config`'s `LogServer`/`NotifServer` capabilities to `flexible-logger`'s engine requirements.
*   **Dynamic Leveling**: Implements `SetLevel` on the core `LogEngine` to allow post-initialization level updates.
