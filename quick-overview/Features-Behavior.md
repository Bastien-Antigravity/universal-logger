---
microservice: universal-logger
type: reference
status: active
tags:
- '#service/universal-logger'
- '#type/reference'
- '#state/active'
- '#zone/3-fleet'
- '#ai/ignore'
---

# Features & Behavior: universal-logger

This document tracks the core capabilities, FFI memory ownership policies, and BDD-style behavioral contracts of the `universal-logger` repository.

---

## 🚀 Core Capabilities

*   **Unified Polyglot Engine**: Employs a single compiled Go dynamic library (`libunilog`) to orchestrate configuration and logging across Go, Python, Rust, C++, and Excel/VBA.
*   **Zero-Allocation Config Access**: Eliminates dynamic allocations on VBA configuration paths by piping requests through a global static memory gateway (`UniLog_Config_Get_Safe`).
*   **Leak-Free Python ctypes Bindings**: Protects Python applications from dynamic FFI leaks by keeping raw pointers bound to `c_void_p` and executing manual deallocation via `DistConf_FreeString`.
*   **Defensive FFI Input Cleaning**: Strips trailing/embedded null-bytes (`\x00`) and trims leading/trailing spaces from all incoming string arguments to prevent memory crashes or pointer mismatches.
*   **Ecosystem Telemetry Parity**: Fully supports standard log levels (DEBUG, STREAM, INFO, LOGON, LOGOUT, TRADE, SCHEDULE, REPORT, WARNING, ERROR, CRITICAL) for seamless cross-tier aggregations.
*   **Asynchronous Live Config Updates**: Supports asynchronous configuration updates using registered callbacks, enabling real-time telemetry updates.

---

## 🛡️ FFI Memory Ownership Protocols

Dynamic boundaries across languages are highly vulnerable to garbage collection desynchronizations and memory leaks. The `universal-logger` enforces three strict memory-safety rules:

### Rule 1: The VBA Static Buffer (Zero C-Heap Leaks)
*   **Given**: An Excel workbook executing high-throughput configuration queries (`GetConfig`).
*   **When**: The query goes through the Go bridge interface.
*   **Then**: The Go core copies the string into a global static array `vbaBuffer` under a thread-safe mutex and returns its pointer. Excel reads it without allocating new memory, resulting in **Zero C-heap leaks**.

### Rule 2: Python `c_void_p` Mapping & Manual Deallocation
*   **Given**: A Python facade querying active configurations.
*   **When**: A value is returned from CGO.
*   **Then**: The return type is treated as `c_void_p` to preserve the raw memory address. After decoding the bytes, the facade manually triggers `lib.DistConf_FreeString(res_ptr)` in a `finally` block, freeing the C-heap memory.

### Rule 3: Defensive Null-Byte Guarding
*   **Given**: A native wrapper calling the shared library with dynamic strings.
*   **When**: The string is received by CGO.
*   **Then**: The Go core runs `sanitizeFFIString`, removing spaces, tabs, newlines, and hidden null terminators to ensure string compatibility and prevent pointer truncation.

---

## 📋 Cross-Language Feature Matrix

| Feature | Go Core | Python Facade | Rust Wrapper | C++ RAII | VBA / Excel |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **FFI Input Cleaning** | ✅ Native | ✅ Via CGO | ✅ Via CGO | ✅ Via CGO | ✅ Via CGO |
| **Config Read (Standard)** | ✅ Native | ✅ Manual Free | ✅ FFI | ✅ FFI | ❌ Not Recommended |
| **Config Read (Safe Buffer)**| ✅ Native | ✅ Supported | ✅ Supported | ✅ Supported | ✅ Enforced (Static) |
| **Async Config Updates** | ✅ Native | ✅ Async Iterator | ✅ FFI Callback | ✅ FFI Callback | ❌ Message Pump Only |
| **Log Level Parity** | ✅ 11 Levels | ✅ 11 Levels | ✅ 11 Levels | ✅ 11 Levels | ✅ 11 Levels |
| **Trace Metadata Injection** | ✅ Zero-Alloc | ✅ Dynamic Dict | ✅ Struct | ✅ RAII Context | ✅ Dict Helper |

---

## ⚙️ Behavior Specifications (BDD Contracts)

### Feature: Configuration Lookups

#### BDD-001: Thread-Safe VBA Lookups
```gherkin
Given a concurrent Excel application session
When multiple sheets query the "GetConfig" API simultaneously
Then the queries must resolve via "UniLog_Config_Get_Safe"
And the Go core must serialize access using the internal mutex
And no dynamic C-heap memory allocations should occur.
```

#### BDD-002: Python Memory Isolation
```gherkin
Given an active Python asyncio loop executing 50,000 queries
When "get_config" is invoked
Then the raw dynamic C pointer must be returned as "c_void_p"
And the value must be safely decoded to a Python string
And the pointer address must be explicitly deallocated using "DistConf_FreeString"
And virtual memory RSS growth must remain below 0.5 MB.
```

#### BDD-003: Robust FFI Input Validation
```gherkin
Given a client library invoking "UniLog_Init"
When a profile name contains leading/trailing whitespaces or embedded null-bytes (e.g. "  standalone\x00  ")
Then "sanitizeFFIString" must clean the string to "standalone"
And the Go bootstrap engine must initialize successfully.
```
