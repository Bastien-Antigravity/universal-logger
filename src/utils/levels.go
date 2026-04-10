package utils

import (
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// -------------------------------------------------------------------------

// Level mirrors the universal-logger Level using a type alias.
type Level = interfaces.Level

// -------------------------------------------------------------------------

// Log level constants mirrored from universal-logger interfaces for easy access.
const (
	LevelNotSet   = interfaces.LevelNotSet
	LevelDebug    = interfaces.LevelDebug
	LevelStream   = interfaces.LevelStream
	LevelInfo     = interfaces.LevelInfo
	LevelLogon    = interfaces.LevelLogon
	LevelLogout   = interfaces.LevelLogout
	LevelTrade    = interfaces.LevelTrade
	LevelSchedule = interfaces.LevelSchedule
	LevelReport   = interfaces.LevelReport
	LevelWarning  = interfaces.LevelWarning
	LevelError    = interfaces.LevelError
	LevelCritical = interfaces.LevelCritical
)

// -------------------------------------------------------------------------

// GetLogLevel converts string to Level.
func GetLogLevel(LogLevel string) Level {
	return logger_models.ParseLevel(LogLevel)
}
