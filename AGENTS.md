---
microservice: universal-logger
type: rules
status: active
tags:
- '#service/universal-logger'
- '#domain/observability'
- '#layer/library'
- '#type/rules'
- '#state/active'
- '#zone/3-fleet'
---

# AGENTS.md: universal-logger

## Service Mission & Architecture Role
`universal-logger` is the standard logging library and abstraction layer across the Bastien-Antigravity ecosystem. It decouples microservices from concrete logging backends by enforcing the unified `ILogger` interface (`Debug`, `Info`, `Warning`, `Error`, `Critical`), supporting structured metadata, multi-sink dispatch (Console, File, Network/SafeSocket to `log-server`), and polyglot CGO bindings (`libunilog`).

- **Ecosystem Role**: Universal logging engine and standard interface.
- **Sinks**:
  - Console / Stdout
  - Rolling File Sink
  - Network Sink (streams over `safe-socket` to `log-server` on port `9020`)
- **Shared CGO Core**: `libunilog` exported for Python, Rust, and C++ callers.
- **Configuration Link**: `standalone.yaml -> ../docker-deployment/modes/local/config/native.yaml`

## Key Build & Test Commands
```bash
# Run tests
go test -v ./...

# Build CLI demo / tool
go build -o bin/universal-logger ./cmd/universal-logger

# Build shared CGO engine
make shared-lib
```

## AI Development & Integration Guidelines
1. **Never Panic in Logger**: Logging failures must fail gracefully to stderr/console and never crash the host service.
2. **SafeSocket Network Sink Stability**: Ensure network reconnection backoff does not block worker threads.
3. **Header Ritual**: All source files MUST begin with the Triple-Block header (`ESSENTIAL PROCESS`, `DATA FLOW`, `KEY PARAMETERS`).
4. **Section Dividers**: Use `// -----------------------------------------------------------------------------` between exported methods.
