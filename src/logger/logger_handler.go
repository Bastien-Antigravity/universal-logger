package logger

import (
	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// UniversalLogger wraps the flexible-logger library.
type UniversalLogger struct {
	Logger interfaces.Logger
}

// -------------------------------------------------------------------------

// NewUniversalLogger initializes a new logger service from an existing logger instance.
func NewUniversalLogger(logger interfaces.Logger) *UniversalLogger {
	return &UniversalLogger{
		Logger: logger,
	}
}

// -------------------------------------------------------------------------
// Debug logs a message at Debug level.
func (s *UniversalLogger) Debug(format string, args ...any) {
	s.Logger.Debug(format, args...)
}

// -------------------------------------------------------------------------
// Info logs a message at Info level.
func (s *UniversalLogger) Info(format string, args ...any) {
	s.Logger.Info(format, args...)
}

// -------------------------------------------------------------------------
// Warning logs a message at Warning level.
func (s *UniversalLogger) Warning(format string, args ...any) {
	s.Logger.Warning(format, args...)
}

// -------------------------------------------------------------------------
// Error logs a message at Error level.
func (s *UniversalLogger) Error(format string, args ...any) {
	s.Logger.Error(format, args...)
}

// -------------------------------------------------------------------------
// Critical logs a message at Critical level.
func (s *UniversalLogger) Critical(format string, args ...any) {
	s.Logger.Critical(format, args...)
}

// -------------------------------------------------------------------------
// Extra functions
// -------------------------------------------------------------------------

// Stream logs a message at Stream level.
func (s *UniversalLogger) Stream(format string, args ...any) {
	s.Logger.Log(logger_models.LevelStream, format, args...)
}

// -------------------------------------------------------------------------
// Logon logs a message at Logon level.
func (s *UniversalLogger) Logon(format string, args ...any) {
	s.Logger.Log(logger_models.LevelLogon, format, args...)
}

// -------------------------------------------------------------------------
// Logout logs a message at Logout level.
func (s *UniversalLogger) Logout(format string, args ...any) {
	s.Logger.Log(logger_models.LevelLogout, format, args...)
}

// -------------------------------------------------------------------------
// Trade logs a message at Trade level.
func (s *UniversalLogger) Trade(format string, args ...any) {
	s.Logger.Log(logger_models.LevelTrade, format, args...)
}

// -------------------------------------------------------------------------
// Schedule logs a message at Schedule level.
func (s *UniversalLogger) Schedule(format string, args ...any) {
	s.Logger.Log(logger_models.LevelSchedule, format, args...)
}

// -------------------------------------------------------------------------
// Report logs a message at Report level.
func (s *UniversalLogger) Report(format string, args ...any) {
	s.Logger.Log(logger_models.LevelReport, format, args...)
}

// -------------------------------------------------------------------------
// Utility functions
// -------------------------------------------------------------------------

// -------------------------------------------------------------------------
// SetLevel sets the current log level.
func (s *UniversalLogger) SetLevel(level logger_models.Level) {
	s.Logger.SetLevel(level)
}

// -------------------------------------------------------------------------
// Close closes the underlying logger.
func (s *UniversalLogger) Close() {
	if s.Logger != nil {
		s.Logger.Close()
	}
}
