package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"

	"github.com/Bastien-Antigravity/universal-logger/src/bootstrap"
	"github.com/Bastien-Antigravity/universal-logger/src/config"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"

	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// FacadeSession holds the state for a single library instantiation.
type FacadeSession struct {
	Config   *config.DistConfig
	Logger   interfaces.Logger
	VbaHwnd  uintptr // Windows HWND for Message Pump (VBA only)
	VbaMsgId uint32  // Windows Message ID (VBA only)
}

var (
	facadeMu    sync.Mutex
	facadeStore         = make(map[uintptr]*FacadeSession)
	facadeId    uintptr = 1
)

func main() {}

// -------------------------------------------------------------------------

//export UniLog_Init
func UniLog_Init(configProfile, appName, loggerProfile *C.char, logLevel C.int, useLocalNotifier C.int) uintptr {
	name := C.GoString(appName)
	cfgProf := C.GoString(configProfile)
	logProf := C.GoString(loggerProfile)
	cfg, log := bootstrap.Init(name, cfgProf, logProf, logger_models.Level(logLevel), useLocalNotifier != 0)

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

// -------------------------------------------------------------------------

//export UniLog_Close
func UniLog_Close(handle uintptr) {
	facadeMu.Lock()
	defer facadeMu.Unlock()
	if session, ok := facadeStore[handle]; ok {
		session.Logger.Close()
		delete(facadeStore, handle)
	}
}
