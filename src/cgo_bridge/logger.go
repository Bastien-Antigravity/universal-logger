package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"universal-logger/src/utils"

	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// -------------------------------------------------------------------------

//export UniLog_LogWithMetadata
func UniLog_LogWithMetadata(handle uintptr, level int, msg, file, line, function, module *C.char) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok || session.Logger == nil {
		return
	}

	utils.LogWithMetadata(
		session.Logger.Logger,
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
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if ok && session.Logger != nil {
		session.Logger.SetLevel(logger_models.Level(level))
	}
}
