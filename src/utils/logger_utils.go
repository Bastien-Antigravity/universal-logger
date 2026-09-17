package utils

// =============================================================================
// ESSENTIAL PROCESS: Logging utilities providing polyglot caller metadata dispatch and facade bridging.
//
// DATA FLOW:
//   1. Receives explicit caller stack metadata from polyglot CGO bridge.
//   2. Delegates execution to LogWithCaller preserving filtering, sampling, and alerts.
//   3. Falls back to standard logging if target does not support caller injection.
//
// KEY PARAMETERS:
//   - LogWithMetadata: Injects explicit caller metadata (file, line, function, module).
// =============================================================================

import (
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
)

// -----------------------------------------------------------------------------

// Logger mirrors the universal-logger main interface using a type alias.
// This allows consumers to use the Logger interface without direct dependency on flexible-logger.
type Logger = interfaces.Logger

// -----------------------------------------------------------------------------

// LogWithMetadata allows manual injection of stack metadata.
// It delegates to LogWithCaller to ensure level filtering, sampling, and alert notifications are respected.
func LogWithMetadata(logger Logger, level Level, msg, file, line, function, module string) {
	var target any = logger
	for target != nil {
		if cl, ok := target.(interface {
			LogWithCaller(level Level, msg, file, line, function, module string)
		}); ok {
			cl.LogWithCaller(level, msg, file, line, function, module)
			return
		}
		if unwrapper, ok := target.(interface{ Unwrap() any }); ok {
			target = unwrapper.Unwrap()
		} else {
			break
		}
	}

	// Fallback to standard logging if caller injection is not supported
	logger.Log(level, "%s", msg)
}

// -----------------------------------------------------------------------------

// Log logs a message at a specific level using the provided logger.
func Log(logger Logger, level Level, format string, args ...any) {
	logger.Log(level, format, args...)
}

// -----------------------------------------------------------------------------

// GetUnderlyingLogger is a helper to access the raw interface (maintained for compatibility/utility).
func GetUnderlyingLogger(logger Logger) Logger {
	return logger
}
