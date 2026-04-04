//go:build !windows
// +build !windows

package main

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
import "unsafe"

// UniLog_RegisterVBAWindow is a stub for non-Windows platforms.
//export UniLog_RegisterVBAWindow
func UniLog_RegisterVBAWindow(handle uintptr, hwnd uintptr, msgId uint32) {
	// No-op on macOS/Linux
}

// dispatchConfigurationUpdate handles the logic of routing a configuration 
// update on non-Windows platforms (Standard FFI Callback only).
func dispatchConfigurationUpdate(handle uintptr, cb C.config_update_cb, jsonData string) {
	if cb != nil {
		cStr := C.CString(jsonData)
		defer C.free(unsafe.Pointer(cStr))
		C.call_config_update_cb(cb, cStr)
	}
}
