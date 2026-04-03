# Testing Documentation: Distconf-Flexlog Facade

This document describes the testing strategy and procedures for the `universal-logger` orchestrator.

## Testing Philosophy

Since the underlying libraries (`distributed-config` and `flexible-logger`) are already extensively unit-tested for their internal logic, the facade tests focus exclusively on **Orchestration** and **Synergy**.

Our goal is to verify that the two systems are correctly "bindingd"—meaning parameters are mapped accurately and the subsystems are initialized in the correct order.

## Test Layers

The project uses a two-layer testing approach:

### 1. Integration & Orchestration (`src/facade/facade_test.go`)
These automated tests verify the internal state of the facade after initialization.
-   **Orchestration Logic**: Ensures that both the Config and Logger subsystems are non-nil and functional.
*   **Binding Verification**: Confirms that discovery data (e.g., LogServer IP/Port) from the config is correctly passed to the logger profile.
-   **Parameter Mapping**: Verifies that `MFacadeParams` fields like `LogLevel` are successfully converted and applied to the engine.
-   **Profile Synergy**: Tests the initialization of key profiles (`devel`, `minimal`) to ensure they play well with the `standalone` config.

### 2. Scenario Verification (`cmd/examples/main.go`)
A comprehensive suite of 6 "real-world" scenarios used to verify the facade's behavior in different environments (Production, Development, Testing, Monitoring).

## Running Tests

### Automated Integration Tests
To run the orchestration suite, use the standard Go toolchain:

```bash
go test -v ./src/facade/...
```

### Manual Scenario Testing
To verify an exhaustive list of operational modes, run the example suite:

```bash
# Example: Run the 'Local Development' scenario
go run cmd/examples/main.go 1

# Example: Run the 'Automated Testing' scenario
go run cmd/examples/main.go 4
```

## Operational Scenarios Table

We verify the following synergistic combinations in every release:

| Scenario | Config Profile | Logger Profile | Primary Purpose |
| :--- | :--- | :--- | :--- |
| **1: Local Dev** | `standalone` | `devel` | Rapid local iteration. |
| **2: Production** | `production` | `standard` | Full remote orchestration. |
| **3: High Load** | `production` | `high_perf` | Low-latency async logging. |
| **4: Testing** | `test` | `minimal` | CI/CD optimized. |
| **5: Monitor** | `preprod` | `notif_logger` | Alert-first monitoring. |
| **6: Critical** | `production` | `no_lock` | Mutex-free execution. |

## Writing New Tests

When adding new integration points (e.g., a new profile or a new parameter mapping):
1.  **Use `standalone` or `test` profiles**: Avoid using `production` in automated tests to prevent mandatory network hangs or config server requirements.
2.  **Verify State**: Check that the resulting `IFacade` object has its underlying subsystems cross-linked (e.g., `facade.Config` matches the config passed to the logger).
3.  **Don't Re-test Logic**: If a bug is found in how a log message is written to a file, the fix and the test belong in `flexible-logger`, not here.
