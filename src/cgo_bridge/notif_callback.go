package main

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
	"unsafe"
	"universal-logger/src/utils"
)

// -------------------------------------------------------------------------

// internalNotificationPump drains the Go channel and pumps messages to the C callback.
func internalNotificationPump(notifQueue <-chan *utils.NotifMessage, callback C.UniLogNotifCallback) {
	for msg := range notifQueue {
		if msg == nil {
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

// -------------------------------------------------------------------------

//export UniLog_RegisterNotifCallback
func UniLog_RegisterNotifCallback(handle uintptr, callback C.UniLogNotifCallback) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok || session.Logger == nil {
		return
	}

	queue := session.Logger.GetNotifQueue()
	if queue == nil {
		// Notifier was not enabled during Init
		return
	}

	// Start a background goroutine to pump notifications
	go internalNotificationPump(queue, callback)
}
