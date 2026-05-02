# 🧬 Project DNA: universal-logger

## 🎯 High-Level Intent (BDD)
- **Goal**: High-performance, zero-overhead polyglot logging for the Bastien-Antigravity ecosystem.
- **Architectural Pattern**: **Super-Bridge** (Consolidated Go binary exporting C-API for multiple languages).
- **Behavioral Source of Truth**: [[business-bdd-brain/02-Behavior-Specs/universal-logger]]
- **Spec Gate**: [HARDENED] No implementation without an `approved` spec in the folder above.

## 🛠️ Role Specifics
- **Architect**: 
    - Focus on CGO boundary safety and memory alignment.
    - Ensure a single Go runtime is shared across all facades.
    - Maintain strict backward compatibility for the C-API.
- **QA**: 
    - Must verify parity across all 5 languages (Go, Python, Rust, C++, VBA).
    - Mandatory run of `scripts/parity-audit.py` before any release-ready state.
- **Developer**:
    - Follow polyglot naming conventions (snake_case for Python, camelCase for Go, etc.).

## 🚦 Lifecycle & Versioning
- **Primary Branch**: `develop`
- **Protected Branches**: `main`, `master`
- **Versioning Strategy**: Semantic Versioning (vX.Y.Z).
- **Version Source of Truth**: `VERSION.txt` (Must be synced to `go.mod`, `Cargo.toml`, and `facade.py`).
