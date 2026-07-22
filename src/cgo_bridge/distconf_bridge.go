package main

/*
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

// Define the callback type for C
typedef void (*distconf_update_cb)(uintptr_t handle, const char* json_data);

// Helper to safely execute a C callback from Go
static void call_distconf_update_cb(distconf_update_cb cb, uintptr_t handle, const char* json_data) {
    if (cb != NULL) {
        cb(handle, json_data);
    }
}

// Standardized Error Codes
#define DISTCONF_SUCCESS                0
#define DISTCONF_ERR_GENERIC            1
#define DISTCONF_ERR_INVALID_HANDLE      2
#define DISTCONF_ERR_KEY_NOT_FOUND       3
#define DISTCONF_ERR_VALIDATION_FAILED   4
#define DISTCONF_ERR_NETWORK_FAILURE     5
#define DISTCONF_ERR_DECRYPTION_FAILED   6
#define DISTCONF_ERR_INVALID_INPUT       7

static char* last_error = NULL;
static int last_error_code = 0;

static void set_last_error(int code, const char* err) {
    last_error_code = code;
    if (last_error != NULL) {
        free(last_error);
    }
    if (err == NULL) {
        last_error = NULL;
    } else {
        last_error = strdup(err);
    }
}
*/
import "C"

import (
	"encoding/json"
	"strings"
	"unsafe"

	"github.com/Bastien-Antigravity/distributed-config/src/cgo_bridge"
)

// -------------------------------------------------------------------------
// HELPERS
// -------------------------------------------------------------------------

func mapErrorCode(err error) C.int {
	if err == nil {
		return C.DISTCONF_SUCCESS
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not found"):
		return C.DISTCONF_ERR_KEY_NOT_FOUND
	case strings.Contains(msg, "invalid handle"):
		return C.DISTCONF_ERR_INVALID_HANDLE
	case strings.Contains(msg, "validation failed"):
		return C.DISTCONF_ERR_VALIDATION_FAILED
	case strings.Contains(msg, "network") || strings.Contains(msg, "connection"):
		return C.DISTCONF_ERR_NETWORK_FAILURE
	case strings.Contains(msg, "decryption"):
		return C.DISTCONF_ERR_DECRYPTION_FAILED
	default:
		return C.DISTCONF_ERR_GENERIC
	}
}

func setLastError(code C.int, msg string) {
	if code == C.DISTCONF_SUCCESS {
		C.set_last_error(code, nil)
	} else {
		cStr := C.CString(msg)
		defer C.free(unsafe.Pointer(cStr))
		C.set_last_error(code, cStr)
	}
}

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
	handle := cgo_bridge.New(sanitizeFFIString(profile))
	if handle == 0 {
		setLastError(C.DISTCONF_ERR_GENERIC, "failed to initialize configuration")
	} else {
		setLastError(C.DISTCONF_SUCCESS, "")
	}
	return handle
}

//export DistConf_Close
func DistConf_Close(handle uintptr) {
	cgo_bridge.Close(handle)
	setLastError(C.DISTCONF_SUCCESS, "")
}

//export DistConf_Get
func DistConf_Get(handle uintptr, section, key *C.char) *C.char {
	val := cgo_bridge.Get(handle, C.GoString(section), C.GoString(key))
	if val == "" {
		setLastError(C.DISTCONF_ERR_KEY_NOT_FOUND, "key not found")
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(val)
}

//export DistConf_Set
func DistConf_Set(handle uintptr, section, key, value *C.char) int {
	if err := cgo_bridge.Set(handle, C.GoString(section), C.GoString(key), C.GoString(value)); err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return 0
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return 1
}

//export DistConf_Sync
func DistConf_Sync(handle uintptr) int {
	if err := cgo_bridge.Sync(handle); err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return 0
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return 1
}

//export DistConf_GetAddress
func DistConf_GetAddress(handle uintptr, capability *C.char) *C.char {
	addr, err := cgo_bridge.GetAddress(handle, C.GoString(capability))
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(addr)
}

//export DistConf_GetGRPCAddress
func DistConf_GetGRPCAddress(handle uintptr, capability *C.char) *C.char {
	addr, err := cgo_bridge.GetGRPCAddress(handle, C.GoString(capability))
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(addr)
}

//export DistConf_GetCapability
func DistConf_GetCapability(handle uintptr, capability *C.char) *C.char {
	val, err := cgo_bridge.GetCapability(handle, C.GoString(capability))
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(val)
}

//export DistConf_GetFullConfig
func DistConf_GetFullConfig(handle uintptr) *C.char {
	val, err := cgo_bridge.GetFullConfig(handle)
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(val)
}

//export DistConf_GetLastError
func DistConf_GetLastError() *C.char {
	return C.last_error
}

//export DistConf_GetLastErrorCode
func DistConf_GetLastErrorCode() int {
	return int(C.last_error_code)
}

//export DistConf_Decrypt
func DistConf_Decrypt(handle uintptr, ciphertext *C.char) *C.char {
	decrypted, err := cgo_bridge.Decrypt(C.GoString(ciphertext))
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(decrypted)
}

//export DistConf_ApplyFileOverride
func DistConf_ApplyFileOverride(handle uintptr, filename *C.char) *C.char {
	localJSON, err := cgo_bridge.ApplyFileOverride(handle, C.GoString(filename))
	if err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return nil
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return C.CString(localJSON)
}

//export DistConf_ShareConfig
func DistConf_ShareConfig(handle uintptr, jsonData *C.char) int {
	if err := cgo_bridge.ShareConfig(handle, C.GoString(jsonData)); err != nil {
		setLastError(mapErrorCode(err), err.Error())
		return 0
	}
	setLastError(C.DISTCONF_SUCCESS, "")
	return 1
}

//export DistConf_OnLiveConfUpdate
func DistConf_OnLiveConfUpdate(handle uintptr, cb C.distconf_update_cb) {
	cgo_bridge.FacadeMu.Lock()
	session, ok := cgo_bridge.FacadeStore[handle]
	cgo_bridge.FacadeMu.Unlock()

	if !ok || session.Config == nil {
		setLastError(C.DISTCONF_ERR_INVALID_HANDLE, "invalid handle")
		return
	}

	setLastError(C.DISTCONF_SUCCESS, "")
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

//export DistConf_OnRegistryUpdate
func DistConf_OnRegistryUpdate(handle uintptr, cb C.distconf_update_cb) {
	cgo_bridge.FacadeMu.Lock()
	session, ok := cgo_bridge.FacadeStore[handle]
	cgo_bridge.FacadeMu.Unlock()

	if !ok || session.Config == nil {
		setLastError(C.DISTCONF_ERR_INVALID_HANDLE, "invalid handle")
		return
	}

	setLastError(C.DISTCONF_SUCCESS, "")
	session.Config.OnRegistryUpdate(func(registry map[string][]string) {
		jsonData, err := json.Marshal(registry)
		if err != nil {
			return
		}

		cStr := C.CString(string(jsonData))
		C.call_distconf_update_cb(cb, C.uintptr_t(handle), cStr)
		C.free(unsafe.Pointer(cStr))
	})
}
