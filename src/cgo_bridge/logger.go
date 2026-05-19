package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// -------------------------------------------------------------------------

//export UniLog_LogWithMetadata
func UniLog_LogWithMetadata(handle uintptr, level int, msg, file, line, function, module *C.char) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if !ok || session.Logger == nil {
		return
	}

	utils.LogWithMetadata(
		session.Logger,
		logger_models.Level(level),
		C.GoString(msg),
		C.GoString(file),
		C.GoString(line),
		C.GoString(function),
		C.GoString(module),
	)
}

// -------------------------------------------------------------------------

//export UniLog_SetLevel
func UniLog_SetLevel(handle uintptr, level int) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if ok && session.Logger != nil {
		session.Logger.SetLevel(logger_models.Level(level))
	}
}

// -------------------------------------------------------------------------

//export UniLog_GetLevel
func UniLog_GetLevel(handle uintptr) int {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if ok && session.Logger != nil {
		return int(session.Logger.GetLevel())
	}
	return 0
}

// -------------------------------------------------------------------------

//export UniLog_AddMetadata
func UniLog_AddMetadata(handle uintptr, key, value *C.char) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if ok && session.Logger != nil {
		session.Logger.AddMetadata(C.GoString(key), C.GoString(value))
	}
}

// -------------------------------------------------------------------------

//export UniLog_SetMetadata
func UniLog_SetMetadata(handle uintptr, jsonMetadata *C.char) {
	facadeMu.RLock()
	session, ok := facadeStore[handle]
	facadeMu.RUnlock()

	if ok && session.Logger != nil {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(C.GoString(jsonMetadata)), &metadata); err == nil {
			session.Logger.SetMetadata(metadata)
		}
	}
}
