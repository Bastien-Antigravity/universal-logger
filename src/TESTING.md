# Testing: Go Core and CGO Bridge

This document explains how to run the Go-side test suite for Universal Logger.

## Automated Unit and Integration Tests

Our primary test suite is contained within the `src/facade/` directory. These tests verify the orchestration between `distributed-config` and `flexible-logger`.

```bash
# Run all Go tests
go test -v ./src/...
```

## Internal Bridge Testing

Because CGO is sensitive to handle management, we use specific tests to ensure:
- **Handle Uniqueness**: No two sessions share the same `uintptr`.
- **Concurrency Safety**: Map access is correctly protected by `sync.Mutex`.
- **Dynamic Config Updating**: Callbacks are properly registered and triggered in Go.

### Running with Coverage
```bash
go test -cover -v ./src/facade/...
```

## Infrastructure Mocking

To ensure that tests can be run in CI/CD without real server dependencies:
- **Config Profiling**: Tests use a `standalone` config profile that loads from local disk.
- **Mock Handlers**: Testing code sometimes implements mock interfaces to verify that data is correctly "binded" from one subsystem to the other.

## Maintenance and Benchmarks

We maintain Go benchmarks to ensure that the orchestration overhead remains negligible.
```bash
go test -bench=. ./src/logger/...
```

## Cross-Platform Considerations

If a test fails on Go, it will likely fail on ALL language facades. Always ensure that the Go core is passing its suite before troubleshooting errors in Python, Rust, or C++.
