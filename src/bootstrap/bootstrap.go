package bootstrap

import (
	"universal-logger/src/utils"

	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
	"github.com/Bastien-Antigravity/universal-logger/src/config"
	"github.com/Bastien-Antigravity/universal-logger/src/logger"
)

// -------------------------------------------------------------------------
// Initialize initializes both subsystems and returns both directly.
// It also sets up the automatic log-level synchronization.
func Init_UniLogger(Name, ConfigProfile, LoggerProfile, LogLevel string) (*config.DistConfig, *logger.UniversalLogger) {
	// 1. Initialize Config Service
	distConfig := config.NewDistributedConfig(ConfigProfile)

	// 2. Initialize Logger using the selected profile
	var flexLogger interfaces.Logger
	switch LoggerProfile {
	case "standard":
		flexLogger = profiles.NewStandardLogger(Name, distConfig.Config)
	case "devel":
		flexLogger = profiles.NewDevelLogger(Name)
	case "high_perf":
		flexLogger = profiles.NewHighPerfLogger(Name, distConfig.Config)
	case "minimal":
		flexLogger = profiles.NewMinimalLogger(Name)
	case "notif_logger":
		flexLogger = profiles.NewNotifLogger(Name, distConfig.Config)
	case "no_lock":
		flexLogger = profiles.NewNoLockLogger(Name, distConfig.Config)
	default:
		flexLogger = profiles.NewStandardLogger(Name, distConfig.Config)
	}

	// 3. Apply initial Log Level
	flexLogger.SetLevel(utils.GetLogLevel(LogLevel))
	universalLogger := logger.NewUniversalLogger(flexLogger)

	// 4. Register automatic LogLevel update from config
	distConfig.OnMemConfUpdate(func(update map[string]map[string]string) {
		if section, ok := update["logger"]; ok {
			if levelStr, ok := section["level"]; ok {
				universalLogger.SetLevel(logger_models.ParseLevel(levelStr))
			}
		}
	})

	return distConfig, universalLogger
}
