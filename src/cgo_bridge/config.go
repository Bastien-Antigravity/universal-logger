package main

/*
#include <stdlib.h>

// Define the callback type for C
typedef void (*config_update_cb)(const char* json_data);
*/
import "C"

import (
	"encoding/json"
)

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
		session.Config.Set(C.GoString(section), C.GoString(key), C.GoString(value))
	}
}

// -------------------------------------------------------------------------

//export UniLog_OnMemConfUpdate
func UniLog_OnMemConfUpdate(handle uintptr, cb C.config_update_cb) {
	println("!!! Go: UniLog_OnMemConfUpdate called for handle:", handle)
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()
	
	if !ok {
		println("!!! Go: Handle NOT FOUND in facadeStore:", handle)
		return
	}
	session.Config.OnMemConfUpdate(func(update map[string]map[string]string) {
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
