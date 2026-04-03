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

import (
	"encoding/json"
	"sync"
	"unsafe"

	"github.com/Bastien-Antigravity/distconf-flexlog/src/facade"
	facade_models "github.com/Bastien-Antigravity/distconf-flexlog/src/models"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

var (
	facadeMu    sync.Mutex
	facadeStore = make(map[uintptr]*facade.DistconfFlexlogFacade)
	facadeId    uintptr = 1
)

func main() {}

// Exported functions for C-Shared library

//export NewFacade
func NewFacade(configProfile, appName, loggerProfile, logLevel *C.char) uintptr {
	p := facade_models.MFacadeParams{
		ConfigProfile: C.GoString(configProfile),
		AppName:       C.GoString(appName),
		LoggerProfile: C.GoString(loggerProfile),
		LogLevel:      C.GoString(logLevel),
	}

	f := facade.NewDistconfFlexlogFacade(p)

	facadeMu.Lock()
	defer facadeMu.Unlock()
	id := facadeId
	facadeStore[id] = f
	facadeId++
	return id
}

//export FreeFacade
func FreeFacade(handle uintptr) {
	facadeMu.Lock()
	defer facadeMu.Unlock()
	if f, ok := facadeStore[handle]; ok {
		f.Close()
		delete(facadeStore, handle)
	}
}

//export LogWithMetadataC
func LogWithMetadataC(handle uintptr, level int, msg, file, line, function, module *C.char) {
	facadeMu.Lock()
	f, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok {
		return
	}

	f.LogWithMetadata(
		logger_models.Level(level),
		C.GoString(msg),
		C.GoString(file),
		C.GoString(line),
		C.GoString(function),
		C.GoString(module),
	)
}

//export RegisterUpdateCallback
func RegisterUpdateCallback(handle uintptr, cb C.config_update_cb) {
	facadeMu.Lock()
	f, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok {
		return
	}

	f.OnMemConfUpdate(func(update map[string]map[string]string) {
		jsonData, err := json.Marshal(update)
		if err != nil {
			return
		}
		cStr := C.CString(string(jsonData))
		defer C.free(unsafe.Pointer(cStr))
		C.call_config_update_cb(cb, cStr)
	})
}
