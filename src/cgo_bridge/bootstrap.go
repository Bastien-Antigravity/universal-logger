package cgo_bridge

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"
	"universal-logger/src/bootstrap"
	"universal-logger/src/config"
	"universal-logger/src/logger"
)

// FacadeSession holds the state for a single library instantiation.
type FacadeSession struct {
	Config *config.DistConfig
	Logger *logger.UniversalLogger
}

var (
	facadeMu    sync.Mutex
	facadeStore         = make(map[uintptr]*FacadeSession)
	facadeId    uintptr = 1
)

func main() {}

// Exported functions for C-Shared library

//export NewFacade
func NewFacade(configProfile, appName, loggerProfile, logLevel *C.char) uintptr {
	name := C.GoString(appName)
	cfgProf := C.GoString(configProfile)
	logProf := C.GoString(loggerProfile)
	logLevelStr := C.GoString(logLevel)

	cfg, log := bootstrap.Initialize(name, cfgProf, logProf, logLevelStr)

	facadeMu.Lock()
	defer facadeMu.Unlock()

	id := facadeId
	facadeStore[id] = &FacadeSession{
		Config: cfg,
		Logger: log,
	}
	facadeId++
	return id
}

//export FreeFacade
func FreeFacade(handle uintptr) {
	facadeMu.Lock()
	defer facadeMu.Unlock()
	if session, ok := facadeStore[handle]; ok {
		session.Logger.Close()
		delete(facadeStore, handle)
	}
}
