package interfaces

import logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"

// -------------------------------------------------------------------------
// Log Level Definitions
// -------------------------------------------------------------------------

// Level mirrors the flexible-logger Level using a type alias.
type Level = logger_models.Level

// Log level constants mirrored for easy access through the Universal interface.
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
// Notification Definitions
// -------------------------------------------------------------------------

// NotifMessage mirrors the flexible-logger NotifMessage using a type alias.
type NotifMessage = logger_models.NotifMessage
