package main

// =============================================================================
// ESSENTIAL PROCESS: Notification callback registry bridging local alerts from Go channels to C function pointers.
//
// DATA FLOW:
//   1. Listens on UniLog NotifQueue channel for high-severity alerts.
//   2. Serializes notification payloads to JSON and dispatches to registered C callbacks.
//
// KEY PARAMETERS:
//   - UniLog_RegisterNotifCallback: Binds external notification handler across C ABI.
// =============================================================================

/*
#include <stdlib.h>

// Typedef for the C callback function
typedef void (*UniLogNotifCallback)(const char* json_msg);

// Helper to safely execute a C callback from Go
static void call_notif_callback(UniLogNotifCallback cb, const char* json_msg) {
    if (cb != NULL) {
        cb(json_msg);
    }
}
*/
import "C"
import (
	"encoding/json"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
	"unsafe"
)

// -----------------------------------------------------------------------------

// internalNotificationPump drains the Go channel and pumps messages to the registered C callback.
func internalNotificationPump(handle uintptr, notifQueue <-chan *utils.NotifMessage) {
	for msg := range notifQueue {
		if msg == nil {
			continue
		}

		facadeMu.RLock()
		session, ok := facadeStore[handle]
		var cbPtr unsafe.Pointer
		if ok && session != nil {
			cbPtr = session.NotifCallback
		}
		facadeMu.RUnlock()
		if !ok {
			return
		}

		callback := C.UniLogNotifCallback(cbPtr)
		if callback == nil {
			continue
		}

		// 1. Serialize message to JSON
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		// 2. Convert to C string
		cStr := C.CString(string(jsonBytes))

		// 3. Execute the C callback (Sync / Fire-and-Forget)
		C.call_notif_callback(callback, cStr)

		// 4. Free the C string immediately (Sync execution is important here)
		C.free(unsafe.Pointer(cStr))
	}
}

// -----------------------------------------------------------------------------

//export UniLog_RegisterNotifCallback
func UniLog_RegisterNotifCallback(handle uintptr, callback C.UniLogNotifCallback) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	if !ok || session.Logger == nil {
		facadeMu.Unlock()
		return
	}

	queue := session.Logger.GetNotifQueue()
	if queue == nil {
		facadeMu.Unlock()
		return
	}

	session.NotifCallback = unsafe.Pointer(callback)
	shouldStart := !session.NotifPumpStarted
	if shouldStart {
		session.NotifPumpStarted = true
	}
	facadeMu.Unlock()

	if shouldStart {
		go internalNotificationPump(handle, queue)
	}
}
