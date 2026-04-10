package bootstrap

import (
	"github.com/Bastien-Antigravity/universal-logger/src/config"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	"github.com/Bastien-Antigravity/universal-logger/src/logger"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"

	flex_interfaces "github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
)

// -------------------------------------------------------------------------

// Init initializes both subsystems and returns both directly.
// It also sets up the automatic log-level synchronization.
// useLocalNotifier: If true, enables an internal 1024-buffered notification queue.
func Init(Name, ConfigProfile, LoggerProfile string, LogLevel logger_models.Level, useLocalNotifier bool) (*config.DistConfig, interfaces.Logger) {
	// 1. Initialize Config Service
	distConfig := config.NewDistributedConfig(ConfigProfile)

	// 2. Initialize Logger using the selected profile
	var flexLogger flex_interfaces.Logger
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

	// 4. Initialize Local Notifier if requested
	if useLocalNotifier {
		// Create a channel with a buffer of 1024
		notifQueue := make(chan *utils.NotifMessage, 1024)
		unilog.NotifQueue = notifQueue

		// Bind the channel only if the logger profile supports it
		// (This calls the type-asserting helper in unilog)
		unilog.SetLocalNotifQueue(notifQueue)
	}

	// 5. Register automatic LogLevel update from config
	distConfig.OnConfigUpdate(func(update map[string]map[string]string) {
		if section, ok := update["logger"]; ok {
			if levelStr, ok := section["level"]; ok {
				unilog.SetLevel(logger_models.ParseLevel(levelStr))
			}
		}
	})

	return distConfig, unilog
}
