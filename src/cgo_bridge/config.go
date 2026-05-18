package main

/*
#include <stdlib.h>

// Define the callback type for C
typedef void (*config_update_cb)(const char* json_data);
*/
import "C"

import (
	"encoding/json"
	"sync"
	"unsafe"
)

var (
	vbaBufferMu sync.Mutex
	vbaBuffer   [4096]byte
)

// -------------------------------------------------------------------------

//export UniLog_Config_Get_Safe
func UniLog_Config_Get_Safe(handle uintptr, section, key *C.char) *C.char {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok || session.Config == nil {
		return nil
	}

	val := session.Config.Get(C.GoString(section), C.GoString(key))
	if val == "" {
		return nil
	}

	vbaBufferMu.Lock()
	defer vbaBufferMu.Unlock()

	// Safe Copy to Static Buffer (truncate if exceeds size)
	copy(vbaBuffer[:], val)
	length := len(val)
	if length >= len(vbaBuffer) {
		length = len(vbaBuffer) - 1
	}
	vbaBuffer[length] = 0 // Enforce NULL terminator

	return (*C.char)(unsafe.Pointer(&vbaBuffer[0]))
}

// -------------------------------------------------------------------------

//export UniLog_Config_Get
func UniLog_Config_Get(handle uintptr, section, key *C.char) *C.char {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

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
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if ok && session.Config != nil {
		// We ignore the error here as the C boundary is void
		_ = session.Config.SetConfig(C.GoString(section), C.GoString(key), C.GoString(value))
	}
}

// -------------------------------------------------------------------------

//export UniLog_OnConfigUpdate
func UniLog_OnConfigUpdate(handle uintptr, cb C.config_update_cb) {
	println("!!! Go: UniLog_OnConfigUpdate called for handle:", handle)
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok {
		println("!!! Go: Handle NOT FOUND in facadeStore:", handle)
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
