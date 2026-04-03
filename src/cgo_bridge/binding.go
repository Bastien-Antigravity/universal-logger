package cgo_bridge

/*
#include <stdlib.h>

// Define the callback type for C
typedef void (*config_update_cb)(const char* json_data);

// C helper to call the callback
static void call_config_update_cb(config_update_cb cb, const char* json_data) {
    if (cb) {
        cb(json_data);
    }
}
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

// Exported functions for Configuration

//export GetConfigValueC
func GetConfigValueC(handle uintptr, section, key *C.char) *C.char {
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

//export RegisterUpdateCallback
func RegisterUpdateCallback(handle uintptr, cb C.config_update_cb) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok || session.Config == nil {
		return
	}

	session.Config.OnMemConfUpdate(func(update map[string]map[string]string) {
		jsonData, err := json.Marshal(update)
		if err != nil {
			return
		}
		cStr := C.CString(string(jsonData))
		defer C.free(unsafe.Pointer(cStr))
		C.call_config_update_cb(cb, cStr)
	})
}
