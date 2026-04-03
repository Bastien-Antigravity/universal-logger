package facade

import (
	"time"
	distributed_config "github.com/Bastien-Antigravity/distributed-config"
	facade_models "github.com/Bastien-Antigravity/distconf-flexlog/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/engine"
	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/flexible-logger/src/error_handler"
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

// LogWithMetadata allows manual injection of stack metadata (useful for Python wrappers).
func (sf *DistconfFlexlogFacade) LogWithMetadata(level logger_models.Level, msg, file, line, function, module string) {
	// 1. Try to access the underlying LogEngine to get the Sink
	if logEngine, ok := sf.Logger.(*engine.LogEngine); ok {
		// 2. Get an entry from the pool
		e := logger_models.EntryPool.Get().(*logger_models.LogEntry)
		e.Reset()

		// 3. Fill basic fields
		e.Level = level
		e.Message = msg
		e.Timestamp = time.Now().UTC()
		e.LoggerName = logEngine.Name
		e.Hostname = logEngine.Hostname
		e.ServiceName = logEngine.ServiceName

		// 4. Fill stack metadata from Python
		e.Filename = file
		e.LineNumber = line
		e.FunctionName = function
		e.Module = module

		// 5. Write to the sink
		if err := logEngine.Sink.Write(e); err != nil {
			error_handler.ReportInternalError(logEngine.Name, "facade.LogWithMetadata", err, msg)
		}
		return
	}

	// Fallback to standard logging if not a LogEngine
	sf.Logger.Log(level, msg)
}

// OnMemConfUpdate registers a callback for configuration updates.
func (sf *DistconfFlexlogFacade) OnMemConfUpdate(fn func(map[string]map[string]string)) {
	sf.Config.OnMemConfUpdate(fn)
}

// Close closes the underlying subsystems.
func (sf *DistconfFlexlogFacade) Close() {
	if sf.Logger != nil {
		sf.Logger.Close()
	}
}

// NewDistconfFlexlogFacade initializes both subsystems.
func NewDistconfFlexlogFacade(p facade_models.MFacadeParams) *DistconfFlexlogFacade {
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
