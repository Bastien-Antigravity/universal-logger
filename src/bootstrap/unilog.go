package bootstrap

import (
	"universal-logger/src/config"
	"universal-logger/src/logger"

	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
)

// -------------------------------------------------------------------------

// Init initializes both subsystems and returns both directly.
// It also sets up the automatic log-level synchronization.
func Init(Name, ConfigProfile, LoggerProfile string, LogLevel logger_models.Level) (*config.DistConfig, *logger.UniLog) {
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
	flexLogger.SetLevel(LogLevel)
	unilog := logger.NewUniLog(flexLogger)

	// 4. Register automatic LogLevel update from config
	distConfig.OnConfigUpdate(func(update map[string]map[string]string) {
		if section, ok := update["logger"]; ok {
			if levelStr, ok := section["level"]; ok {
				unilog.SetLevel(logger_models.ParseLevel(levelStr))
			}
		}
	})

	return distConfig, unilog
}
