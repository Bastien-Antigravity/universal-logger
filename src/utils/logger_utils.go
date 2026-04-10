package utils

import (
	"time"

	"github.com/Bastien-Antigravity/flexible-logger/src/engine"
	"github.com/Bastien-Antigravity/flexible-logger/src/error_handler"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// -------------------------------------------------------------------------

// Logger mirrors the universal-logger main interface using a type alias.
// This allows consumers to use the Logger interface without direct dependency on flexible-logger.
type Logger = interfaces.Logger

// -------------------------------------------------------------------------

// LogWithMetadata allows manual injection of stack metadata.
// It tries to access the underlying LogEngine sink for high-performance writing.
func LogWithMetadata(logger Logger, level Level, msg, file, line, function, module string) {
	// 1. Try to access the underlying LogEngine to get the Sink
	if logEngine, ok := logger.(*engine.LogEngine); ok {
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

		// 4. Fill stack metadata
		e.Filename = file
		e.LineNumber = line
		e.FunctionName = function
		e.Module = module

		// 5. Write to the sink
		if err := logEngine.Sink.Write(e); err != nil {
			error_handler.ReportInternalError(logEngine.Name, "logger_utils.LogWithMetadata", err, msg)
		}
		return
	}

	// Fallback to standard logging if not a LogEngine
	logger.Log(level, msg)
}

// -------------------------------------------------------------------------

// Log logs a message at a specific level using the provided logger.
func Log(logger Logger, level Level, format string, args ...any) {
	logger.Log(level, format, args...)
}


// -------------------------------------------------------------------------

// GetUnderlyingLogger is a helper to access the raw interface (maintained for compatibility/utility).
func GetUnderlyingLogger(logger Logger) Logger {
	return logger
}

