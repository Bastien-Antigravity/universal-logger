package main

// =============================================================================
// ESSENTIAL PROCESS: CGO export functions for distributed configuration access and dynamic update callbacks.
//
// DATA FLOW:
//   1. Queries and updates configuration keys on behalf of polyglot consumers.
//   2. Bridges Go configuration update channels to C function pointer callbacks.
//
// KEY PARAMETERS:
//   - UniLog_Config_Get: Fetches configuration values across C ABI.
//   - UniLog_Config_Set: Updates configuration values in memory.
//   - UniLog_OnConfigUpdate: Registers update notification callbacks.
// =============================================================================


/*
#include <stdlib.h>
#include <string.h>

#if defined(__STDC_VERSION__) && __STDC_VERSION__ >= 201112L
#define THREAD_LOCAL _Thread_local
#elif defined(_MSC_VER)
#define THREAD_LOCAL __declspec(thread)
#elif defined(__GNUC__)
#define THREAD_LOCAL __thread
#else
#define THREAD_LOCAL
#endif

static THREAD_LOCAL char vba_buffer[4096];

static inline char* get_vba_buffer(void) {
    return vba_buffer;
}

// Define the callback type for C
typedef void (*config_update_cb)(const char* json_data);
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

// -------------------------------------------------------------------------

//export UniLog_Config_Get_Safe
func UniLog_Config_Get_Safe(handle uintptr, section, key *C.char) *C.char {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if !ok || session.Config == nil {
		return nil
	}

	val := session.Config.Get(C.GoString(section), C.GoString(key))
	if val == "" {
		return nil
	}

	bufPtr := C.get_vba_buffer()
	b := []byte(val)
	if len(b) >= 4096 {
		b = b[:4095]
	}
	C.memcpy(unsafe.Pointer(bufPtr), unsafe.Pointer(&b[0]), C.size_t(len(b)))
	*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(bufPtr)) + uintptr(len(b)))) = 0

	return bufPtr
}

// -------------------------------------------------------------------------

//export UniLog_Config_Get
func UniLog_Config_Get(handle uintptr, section, key *C.char) *C.char {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if !ok || session.Config == nil {
		return nil
	}

	val := session.Config.Get(C.GoString(section), C.GoString(key))
	if val == "" {
		return nil
	}
	return C.CString(val)
}

// -------------------------------------------------------------------------

//export UniLog_Config_Set
func UniLog_Config_Set(handle uintptr, section, key, value *C.char) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if ok && session.Config != nil {
		// We ignore the error here as the C boundary is void
		_ = session.Config.SetConfig(C.GoString(section), C.GoString(key), C.GoString(value))
	}
}

// -------------------------------------------------------------------------

//export UniLog_OnConfigUpdate
func UniLog_OnConfigUpdate(handle uintptr, cb C.config_update_cb) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if !ok || session.Config == nil {
		return
	}
	session.Config.OnConfigUpdate(func(update map[string]map[string]string) {
		jsonData, err := json.Marshal(update)
		if err != nil {
			return
		}

		// Run callback in a goroutine to avoid deadlocks with the Python GIL
		go func() {
			// Delegate all dispatching (FFI + VBA) to the unified dispatcher
			// This ensures config.go stays clean and unaware of VBA internals.
			dispatchConfigurationUpdate(handle, cb, string(jsonData))
		}()
	})
}
