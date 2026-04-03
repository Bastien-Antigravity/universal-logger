package cgo_bridge

/*
#include <stdlib.h>
*/
import "C"

import (
	"universal-logger/src/utils"
)

// Exported functions for Logging

//export LogWithMetadataC
func LogWithMetadataC(handle uintptr, level int, msg, file, line, function, module *C.char) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if !ok || session.Logger == nil {
		return
	}

	utils.LogWithMetadata(
		session.Logger.Logger,
		utils.GetLogLevel(level),
		C.GoString(msg),
		C.GoString(file),
		C.GoString(line),
		C.GoString(function),
		C.GoString(module),
	)
}

//export SetLevelC
func SetLevelC(handle uintptr, level int) {
	facadeMu.Lock()
	session, ok := facadeStore[handle]
	facadeMu.Unlock()

	if ok && session.Logger != nil {
		session.Logger.SetLevel(utils.GetLogLevel(level))
	}
}
