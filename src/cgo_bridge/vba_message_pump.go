//go:build windows
// +build windows

package main

/*
#include <windows.h>
#include <stdlib.h>

// Helper to call PostMessageA from CGO on Windows
static void PostVbaMessage(HWND hwnd, UINT msg, const char* json_data) {
    if (hwnd) {
        PostMessageA(hwnd, msg, 0, (LPARAM)json_data);
    }
}

// C helper to call the configuration update callback
static void call_config_update_cb(config_update_cb cb, const char* json_data) {
    if (cb) {
        cb(json_data);
    }
}
*/
import "C"
import (
	"unsafe"
)

// UniLog_RegisterVBAWindow registers a Windows HWND and Message ID for 
// receiving asynchronous configuration updates in VBA.
//export UniLog_RegisterVBAWindow
func UniLog_RegisterVBAWindow(handle uintptr, hwnd uintptr, msgId uint32) {
	facadeMu.Lock()
	defer facadeMu.Unlock()

	if session, ok := facadeStore[handle]; ok {
		println("!!! Go: Registering VBA Message Pump for handle:", handle, "HWND:", hwnd)
		session.VbaHwnd = hwnd
		session.VbaMsgId = msgId
	}
}

// dispatchConfigurationUpdate handles the logic of routing a configuration 
// update to the correct destination (VBA Message Pump and/or Standard FFI Callback).
func dispatchConfigurationUpdate(handle uintptr, cb C.config_update_cb, jsonData string) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok {
		return
	}

	// 1. Dispatch to VBA (if registered)
	if session.VbaHwnd != 0 {
		// Create a C-string for the JSON data.
		// Note: The VBA side is responsible for consuming this.
		cStr := C.CString(jsonData)

		// Since PostMessage is asynchronous, we do NOT free cStr here.
		// The VBA side receives the pointer in its message loop.
		C.PostVbaMessage(C.HWND(unsafe.Pointer(session.VbaHwnd)), C.UINT(session.VbaMsgId), cStr)
	}

	// 2. Dispatch to standard FFI callback (if registered)
	if cb != nil {
		cStr := C.CString(jsonData)
		defer C.free(unsafe.Pointer(cStr))
		C.call_config_update_cb(cb, cStr)
	}
}
