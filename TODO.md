# TODO: universal-logger

## 🚨 High Priority (Governance Gaps)
- [ ] **VBA Buffer (Purger Rule)**: Replace dynamic C.CString with a fixed-length shared buffer in the Go bridge to eliminate memory leaks without requiring VBA-side freeing (FEAT-003). (Approval Required)
- [ ] **FFI Sanitization**: Add strings.TrimSpace and NULL-byte cleaning to all incoming profile names (FEAT-002). (Approval Required)
- [ ] **Python Memory Leak**: Audit the `unilog.py` wrapper and ensure EVERY call to Get/GetAddress is followed by a `FreeString` call (FEAT-011). (Approval Required)

## 🏗️ Architecture & Refactoring
- [ ] Implement C++ RAII wrapper for session management.
- [ ] Finalize Rust `unilog-rs` crate documentation.

## 🧪 Testing & CI/CD
- [ ] Create cross-language build-pipeline for DLL/SO generation.

## ✅ Completed
- [x] Initial BDD Spec migration.
