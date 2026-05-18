package main

/*
#include <stdlib.h>
#include <stdint.h>

// Define the callback type for C
typedef void (*distconf_update_cb)(uintptr_t handle, const char* json_data);

// Helper to safely execute a C callback from Go
static void call_distconf_update_cb(distconf_update_cb cb, uintptr_t handle, const char* json_data) {
    if (cb != NULL) {
        cb(handle, json_data);
    }
}
*/
import "C"

import (
	"encoding/json"
	"github.com/Bastien-Antigravity/distributed-config/src/cgo_bridge"
	"unsafe"
)

// -------------------------------------------------------------------------
// DISTRIBUTED CONFIG BRIDGE (Re-exported logic from distributed-config)
// -------------------------------------------------------------------------

//export DistConf_FreeString
func DistConf_FreeString(ptr *C.char) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

//export DistConf_New
func DistConf_New(profile *C.char) uintptr {
	return cgo_bridge.New(sanitizeFFIString(profile))
}

//export DistConf_Close
func DistConf_Close(handle uintptr) {
	cgo_bridge.Close(handle)
}

//export DistConf_Get
func DistConf_Get(handle uintptr, section, key *C.char) *C.char {
	val := cgo_bridge.Get(handle, C.GoString(section), C.GoString(key))
	if val == "" {
		return nil
	}
	return C.CString(val)
}

//export DistConf_Set
func DistConf_Set(handle uintptr, section, key, value *C.char) int {
	if err := cgo_bridge.Set(handle, C.GoString(section), C.GoString(key), C.GoString(value)); err != nil {
		return 0
	}
	return 1
}

//export DistConf_Sync
func DistConf_Sync(handle uintptr) int {
	if err := cgo_bridge.Sync(handle); err != nil {
		return 0
	}
	return 1
}

//export DistConf_GetAddress
func DistConf_GetAddress(handle uintptr, capability *C.char) *C.char {
	addr, err := cgo_bridge.GetAddress(handle, C.GoString(capability))
	if err != nil {
		return nil
	}
	return C.CString(addr)
}

//export DistConf_GetFullConfig
func DistConf_GetFullConfig(handle uintptr) *C.char {
	val, err := cgo_bridge.GetFullConfig(handle)
	if err != nil {
		return nil
	}
	return C.CString(val)
}

//export DistConf_Decrypt
func DistConf_Decrypt(handle uintptr, ciphertext *C.char) *C.char {
	decrypted, err := cgo_bridge.Decrypt(C.GoString(ciphertext))
	if err != nil {
		return nil
	}
	return C.CString(decrypted)
}

//export DistConf_OnLiveConfUpdate
func DistConf_OnLiveConfUpdate(handle uintptr, cb C.distconf_update_cb) {
	// Note: callbacks require more complex handling if the logic is in a separate package,
	// because the C callback function pointer needs to be passed correctly.
	// However, since we are building one binary, we can define the listener here.

	cgo_bridge.FacadeMu.Lock()
	session, ok := cgo_bridge.FacadeStore[handle]
	cgo_bridge.FacadeMu.Unlock()

	if !ok || session.Config == nil {
		return
	}

	session.Config.OnLiveConfUpdate(func(update map[string]map[string]string) {
		jsonData, err := json.Marshal(update)
		if err != nil {
			return
		}

		cStr := C.CString(string(jsonData))
		C.call_distconf_update_cb(cb, C.uintptr_t(handle), cStr)
		C.free(unsafe.Pointer(cStr))
	})
}
