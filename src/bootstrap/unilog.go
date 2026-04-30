package bootstrap

import (
	"github.com/Bastien-Antigravity/universal-logger/src/config"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	"github.com/Bastien-Antigravity/universal-logger/src/logger"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"

	flexible_logger "github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// BootstrapOptions defines the configuration parameters for initializes the subsystems.
type BootstrapOptions struct {
	Name             string
	ConfigProfile    string
	LoggerProfile    string
	InitialLogLevel  interfaces.Level
	UseLocalNotifier bool
	ExistingConfig   *config.DistConfig // OPTIONAL: Inject an existing configuration instance
	Metadata         map[string]string  // OPTIONAL: Fields to be added to all logs
}

// -----------------------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------------------

// Init initializes both subsystems and returns both directly.
// It also sets up the automatic log-level synchronization.
// useLocalNotifier: If true, enables an internal 1024-buffered notification queue.
// existingConfig: OPTIONAL. If provided, the logger will use this configuration instance instead of creating a new one.
func Init(Name, ConfigProfile, LoggerProfile, LogLevel string, useLocalNotifier bool, existingConfig *config.DistConfig) (*config.DistConfig, interfaces.Logger) {
	return InitWithOptions(BootstrapOptions{
		Name:             Name,
		ConfigProfile:    ConfigProfile,
		LoggerProfile:    LoggerProfile,
		InitialLogLevel:  utils.GetLogLevel(LogLevel),
		UseLocalNotifier: useLocalNotifier,
		ExistingConfig:   existingConfig,
	})
}

// InitWithOptions initializes the subsystems using a sophisticated options structure.
// It supports dependency injection (ExistingConfig) and specialized metadata.
func InitWithOptions(opts BootstrapOptions) (*config.DistConfig, interfaces.Logger) {
	// 1. Initialize or Inject Config Service
	var distConfig *config.DistConfig
	if opts.ExistingConfig != nil {
		distConfig = opts.ExistingConfig
	} else {
		distConfig = config.NewDistributedConfig(opts.ConfigProfile)
	}

	// 2. Synchronize Service Name
	if opts.Name != "" && distConfig.Config != nil {
		distConfig.Config.Common.Name = opts.Name
	}

	// 3. Initialize Logger using the selected profile
	var flexLogger flexible_logger.Logger
	switch opts.LoggerProfile {
	case "audit":
		flexLogger = profiles.NewAuditLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	case "cloud", "cloud_native":
		flexLogger = profiles.NewCloudLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	case "devel":
		flexLogger = profiles.NewDevelLogger(opts.Name, opts.UseLocalNotifier)
	case "high_perf":
		flexLogger = profiles.NewHighPerfLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	case "minimal":
		flexLogger = profiles.NewMinimalLogger(opts.Name, opts.UseLocalNotifier)
	case "no_lock":
		flexLogger = profiles.NewNoLockLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	case "notif_logger":
		flexLogger = profiles.NewNotifLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	case "standard":
		flexLogger = profiles.NewStandardLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	default:
		flexLogger = profiles.NewStandardLogger(opts.Name, distConfig.Config, opts.UseLocalNotifier)
	}

	// 4. Apply initial Log Level and Metadata
	flexLogger.SetLevel(opts.InitialLogLevel)
	unilog := logger.NewUniLog(flexLogger)

	if len(opts.Metadata) > 0 {
		unilog.SetMetadata(opts.Metadata)
	}

	// 5. Initialize Local Notifier if requested
	if opts.UseLocalNotifier {
		// Create a channel with a buffer of 1024
		notifQueue := make(chan *utils.NotifMessage, 1024)
		unilog.NotifQueue = notifQueue

		// Bind the channel only if the logger profile supports it
		unilog.SetLocalNotifQueue(notifQueue)
	}

	// 6. Register automatic LogLevel update from config
	distConfig.OnConfigUpdate(func(update map[string]map[string]string) {
		if section, ok := update["logger"]; ok {
			if levelStr, ok := section["level"]; ok {
				unilog.SetLevel(logger_models.ParseLevel(levelStr))
			}
		}
	})

	return distConfig, unilog
}
