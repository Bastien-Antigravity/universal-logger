package facade

import (
	distributed_config "github.com/Bastien-Antigravity/distributed-config"
	"github.com/Bastien-Antigravity/distconf-flexlog/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
)

// DistconfFlexlogFacade orchestrates both configuration and logging.
type DistconfFlexlogFacade struct {
	interfaces.Logger
	Config *distributed_config.Config
}

// GetConfig returns the underlying distributed configuration.
func (sf *DistconfFlexlogFacade) GetConfig() *distributed_config.Config {
	return sf.Config
}

// NewDistconfFlexlogFacade initializes both subsystems.
func NewDistconfFlexlogFacade(p models.MFacadeParams) *DistconfFlexlogFacade {
	// 1. Initialize Distributed Config
	cfg := distributed_config.New(p.ConfigProfile)

	// 2. Initialize Logger using the selected profile
	var logger interfaces.Logger
	switch p.LoggerProfile {
	case "standard":
		logger = profiles.NewStandardLogger(p.AppName, cfg)
	case "devel":
		logger = profiles.NewDevelLogger(p.AppName)
	case "high_perf":
		logger = profiles.NewHighPerfLogger(p.AppName, cfg)
	case "minimal":
		logger = profiles.NewMinimalLogger(p.AppName)
	case "notif_logger":
		logger = profiles.NewNotifLogger(p.AppName, cfg)
	case "no_lock":
		logger = profiles.NewNoLockLogger(p.AppName, cfg)
	default:
		// Default to standard if unknown
		logger = profiles.NewStandardLogger(p.AppName, cfg)
	}

	// 3. Apply custom Log Level
	logger.SetLevel(p.GetLogLevel())

	return &DistconfFlexlogFacade{
		Logger: logger,
		Config: cfg,
	}
}
