---
microservice: universal-logger
type: task
status: active
tags:
- '#service/universal-logger'
- '#type/task'
- '#state/active'
- '#zone/3-fleet'
---

# TODO: universal-logger

## 🚨 High Priority (Governance Gaps)
- [x] **VBA Buffer (Purger Rule)**: Replace dynamic C.CString with a fixed-length shared buffer in the Go bridge to eliminate memory leaks without requiring VBA-side freeing (FEAT-003).
- [x] **FFI Sanitization**: Add strings.TrimSpace and NULL-byte cleaning to all incoming profile names (FEAT-002).
- [x] **Python Memory Leak**: Audit the `unilog.py` wrapper and ensure EVERY call to Get/GetAddress is followed by a `FreeString` call (FEAT-011).

## 🏗️ Architecture & Refactoring
- [x] Implement C++ RAII wrapper for session management (FEAT-012).
- [x] Finalize Rust `unilog-rs` crate documentation (FEAT-013).

## 🧪 Testing & CI/CD
- [x] Create cross-language build-pipeline for DLL/SO/dylib generation via consolidated Makefile (FEAT-014).

## ✅ Completed
- [x] **Go Core & CGO Bridge Hardening**: Fixed `UniLog` caller stack skip (`SetCallerSkip(2)` + unified `logInternal`), added lock-free `atomic.Pointer` pre-computed metadata snapshots and hot-path level gates (`UniLog`, `UniLog_LogWithMetadata`, Python `_log`/`_async_log`), unbound local notification queues safely in `UniLog.Close()` without GC finalizers, guarded `vba_buffer` and `last_error_buf` with `sync.Mutex`, prevented duplicate `internalNotificationPump` goroutines, and gated `OwnsConfig` on `UniLog_Close`.
- [x] **Polyglot Parity & Compliance**: Added `#[track_caller]` + all 11 level helpers/macros to Rust `unilog-rs` and C++ `UniversalLogger.hpp`, standardized Python granular camelCase imports and 95-dash dividers, resolved `setup.py` `HARDCODED_VERSION_TRAVERSAL`, fixed `scripts/build_all.sh` and `scripts/parity-audit.py`, removed subpackage `standalone.yaml` symlinks and duplicate stub docs, and isolated Go/Python/Rust/C++ tests via temp `CONFIG_PATH` / `LOG_DIR`.
- [x] Initial BDD Spec migration.
- [x] **VBA Buffer (Purger Rule)**: Replace dynamic C.CString with a fixed-length shared buffer in the Go bridge (FEAT-003).
- [x] **FFI Sanitization**: Add strings.TrimSpace and NULL-byte cleaning to all incoming profile names (FEAT-002).
- [x] **Python Memory Leak**: Audit Python wrapper and decode/free raw void pointers to prevent heap leaks (FEAT-011).
- [x] C++ RAII wrapper for dynamic logger sessions under `unilog/cpp/`.
- [x] Rust `unilog-rs` crate FFI build linking pop-logic under `unilog/rust/`.
- [x] Cross-platform Makefile dynamic OS library generation and dynamic `install_name_tool` dynamic ID registration.
- [x] Reorganized all facades under consolidated `unilog/` layout.
