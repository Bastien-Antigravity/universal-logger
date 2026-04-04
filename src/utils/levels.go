package utils

import logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"

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

// GetLogLevel converts string to logger_models.Level.
func GetLogLevel(LogLevel string) logger_models.Level {
	return logger_models.ParseLevel(LogLevel)
}
