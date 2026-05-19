package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"sync"

	"github.com/Bastien-Antigravity/universal-logger/src/bootstrap"
	"github.com/Bastien-Antigravity/universal-logger/src/config"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"

	distconf_bridge "github.com/Bastien-Antigravity/distributed-config/src/cgo_bridge"
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
	facadeMu    sync.RWMutex
	facadeStore         = make(map[uintptr]*FacadeSession)
	// Start at a high base to avoid collisions with distributed-config IDs (which start at 1)
	facadeId uintptr = 0x10000
)

func main() {}

// -------------------------------------------------------------------------

// sanitizeFFIString cleans C strings coming across the frontier
func sanitizeFFIString(cStr *C.char) string {
	if cStr == nil {
		return ""
	}
	goStr := C.GoString(cStr)
	// Remove standard null chars and embedded controls
	goStr = strings.ReplaceAll(goStr, "\x00", "")
	// Trim surrounding spaces, tabs, and newlines
	return strings.TrimSpace(goStr)
}

//export UniLog_Init
func UniLog_Init(appName, configProfile, loggerProfile *C.char, logLevel C.int, useLocalNotifier C.int, configHandle uintptr) uintptr {
	name := sanitizeFFIString(appName)
	cfgProf := sanitizeFFIString(configProfile)
	logProf := sanitizeFFIString(loggerProfile)

	// Attempt to recover an existing config if handle is provided
	var existingCfg *config.DistConfig
	if configHandle != 0 {
		// 1. First, check if it's one of OUR handles (from UniLog_Init or UniLog_Config_New)
		facadeMu.RLock()
		if session, ok := facadeStore[configHandle]; ok && session.Config != nil {
			existingCfg = session.Config
		}
		facadeMu.RUnlock()

		// 2. If not found, check if it's a handle from the distributed-config bridge
		if existingCfg == nil {
			distconf_bridge.FacadeMu.RLock()
			if session, ok := distconf_bridge.FacadeStore[configHandle]; ok && session.Config != nil {
				existingCfg = &config.DistConfig{Config: session.Config}
			}
			distconf_bridge.FacadeMu.RUnlock()
		}

		// SAFETY: NEVER attempt a raw pointer cast here. Handles are IDs, not memory addresses.
	}

	cfg, log := bootstrap.Init(name, cfgProf, logProf, logger_models.Level(logLevel).String(), useLocalNotifier != 0, existingCfg)

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
