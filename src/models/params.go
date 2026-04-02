package models

import "github.com/Bastien-Antigravity/flexible-logger/src/models"

// MFacadeParams defines the initialization parameters for the facade.
type MFacadeParams struct {
	// Distributed Config
	ConfigProfile string // standalone | test | preprod | production
	
	// Logger
	AppName       string // Application identifier used by both systems
	LoggerProfile string // standard | devel | high_perf | minimal | notif_logger | no_lock
	LogLevel      string // debug | info | warning | error | critical
	
	// Advanced
	PublicIP string // Optional: used for remote identification
}

// GetLogLevel converts the string representation to the internal models.Level
func (m *MFacadeParams) GetLogLevel() models.Level {
	return models.ParseLevel(m.LogLevel)
}
