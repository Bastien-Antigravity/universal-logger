package utils

import logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"

// -------------------------------------------------------------------------

// Level mirrors the flexible-logger Level using a type alias.
// This allowing consumers to use the Level type without direct dependency on flexible-logger.
type Level = logger_models.Level

// -------------------------------------------------------------------------

// Log level constants mirrored from flexible-logger for easy access.
const (
	LevelNotSet   = logger_models.LevelNotSet
	LevelDebug    = logger_models.LevelDebug
	LevelStream   = logger_models.LevelStream
	LevelInfo     = logger_models.LevelInfo
	LevelLogon    = logger_models.LevelLogon
	LevelLogout   = logger_models.LevelLogout
	LevelTrade    = logger_models.LevelTrade
	LevelSchedule = logger_models.LevelSchedule
	LevelReport   = logger_models.LevelReport
	LevelWarning  = logger_models.LevelWarning
	LevelError    = logger_models.LevelError
	LevelCritical = logger_models.LevelCritical
)

// -------------------------------------------------------------------------

// GetLogLevel converts string to Level.
func GetLogLevel(LogLevel string) Level {
	return logger_models.ParseLevel(LogLevel)
}
