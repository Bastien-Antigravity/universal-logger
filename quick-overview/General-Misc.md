---
microservice: universal-logger
type: general-misc
status: active
tags:
- '#service/universal-logger'
- '#type/general-misc'
- '#state/active'
- '#zone/3-fleet'
- '#ai/ignore'
---

# General Miscellaneous: universal-logger

This document contains fallback mechanics, diagnostic logging routines, and frequently asked questions for the `universal-logger`.

---

## 🚨 Internal Diagnostics & Fallback Routines

When writing a dynamic logging facade that spans multiple language wrappers, preventing recursive error loops is critical. For instance, if an FFI bridge logging call encounters a networking error, trying to log that network failure using the same bridge would trigger infinite recursion, leading to thread exhaustion or stack overflows.

### 1. Fallback Path
*   **Mechanism**: If a dynamic CGO boundary call encounters a severe internal system failure (such as mutex blockages, bad profile mapping, or shared library load errors), the system automatically redirects the diagnostic alert away from standard network/file sinks.
*   **Destination**: Written directly to **`os.Stderr`** (Standard Error).
*   **Format**: Prefixed with `[FFI-BRIDGE-INTERNAL-ERROR]` as unformatted raw string output, guaranteeing delivery without utilizing additional dynamic heap space.

### 2. Operational Impact
*   If your client applications (e.g. Python asyncio loop or VBA Excel worksheet) behave unexpectedly during profile updates, monitor the standard console output. Internal FFI errors are explicitly printed to standard error to prevent silent data loss while keeping the primary thread non-blocking.

---

## 🏷️ Ecosystem Tag Taxonomy

Observability logs and configuration schemas across the Bastien-Antigravity fleet rely on standard tags to categorize notes, schemas, and metrics:

*   **Service Identifier Tag**: Always label the service using `#service/universal-logger`.
*   **Domain Categorization**: Organized under the `#domain/observability` umbrella.
*   **Ecosystem Layer**: Placed under `#zone/3-fleet` (Global Fleet CI/CD and Ops).

---

## 🛠️ Frequently Asked Questions (FAQ)

### How does the VBA buffer scale under high concurrent access?
The global static buffer `vbaBuffer` (4KB) is protected by a standard Go `sync.Mutex` inside the compiled dynamic core. If multiple Excel worksheets attempt to query configurations concurrently:
1.  Access to the buffer is strictly serialized.
2.  The value is copied to the buffer and immediately read by VBA using `StringFromPtr`.
3.  The mutex is unlocked, allowing the next query to proceed.
This ensures complete thread-safety with zero risk of concurrency race conditions or memory overlaps.

### Why do we map dynamic string lookups to `c_void_p` in Python?
Setting the ctypes return type to `c_char_p` tells Python to automatically convert the returned character pointer into a Python `str` object. During this automated conversion, Python loses the underlying C-heap memory address, making it impossible to deallocate the memory and causing a persistent leak. Setting the return type to `c_void_p` retains the raw memory address, allowing the facade to decode it using `ctypes.string_at()` and manually free it via `DistConf_FreeString`.

### Can I run the logger in a blocking synchronous mode?
Yes. If you instantiate the **Audit** configuration profile inside `standalone.yaml` or a config server, all logging and notification pipelines operate in a synchronous, blocking mode. In this mode, no messages are dropped, and execution blocks until successful write verification is returned.
