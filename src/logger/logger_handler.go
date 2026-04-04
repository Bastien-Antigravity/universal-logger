package logger

import (
	"runtime"

	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// UniLog wraps the flexible-logger library.
type UniLog struct {
	Logger interfaces.Logger
}

// -------------------------------------------------------------------------

// NewUniLog initializes a new logger service from an existing logger instance.
// Note: This implementation uses a runtime finalizer to automatically call Close() 
// when the logger instance is about to be garbage collected.
func NewUniLog(logger interfaces.Logger) *UniLog {
	res := &UniLog{
		Logger: logger,
	}

	// Register finalizer for automatic cleanup
	runtime.SetFinalizer(res, func(ul *UniLog) {
		ul.Close()
	})

	return res
}


// -------------------------------------------------------------------------

// Debug logs a message at Debug level.
func (s *UniLog) Debug(format string, args ...any) {
	s.Logger.Debug(format, args...)
}

// -------------------------------------------------------------------------

// Info logs a message at Info level.
func (s *UniLog) Info(format string, args ...any) {
	s.Logger.Info(format, args...)
}

// -------------------------------------------------------------------------

// Warning logs a message at Warning level.
func (s *UniLog) Warning(format string, args ...any) {
	s.Logger.Warning(format, args...)
}

// -------------------------------------------------------------------------

// Error logs a message at Error level.
func (s *UniLog) Error(format string, args ...any) {
	s.Logger.Error(format, args...)
}

// -------------------------------------------------------------------------

// Critical logs a message at Critical level.
func (s *UniLog) Critical(format string, args ...any) {
	s.Logger.Critical(format, args...)
}

// -------------------------------------------------------------------------

// Stream logs a message at Stream level.
func (s *UniLog) Stream(format string, args ...any) {
	s.Logger.Log(logger_models.LevelStream, format, args...)
}

// -------------------------------------------------------------------------

// Logon logs a message at Logon level.
func (s *UniLog) Logon(format string, args ...any) {
	s.Logger.Log(logger_models.LevelLogon, format, args...)
}

// -------------------------------------------------------------------------

// Logout logs a message at Logout level.
func (s *UniLog) Logout(format string, args ...any) {
	s.Logger.Log(logger_models.LevelLogout, format, args...)
}

// -------------------------------------------------------------------------

// Trade logs a message at Trade level.
func (s *UniLog) Trade(format string, args ...any) {
	s.Logger.Log(logger_models.LevelTrade, format, args...)
}

// -------------------------------------------------------------------------

// Schedule logs a message at Schedule level.
func (s *UniLog) Schedule(format string, args ...any) {
	s.Logger.Log(logger_models.LevelSchedule, format, args...)
}

// -------------------------------------------------------------------------

// Report logs a message at Report level.
func (s *UniLog) Report(format string, args ...any) {
	s.Logger.Log(logger_models.LevelReport, format, args...)
}

// -------------------------------------------------------------------------
// -------------------------------------------------------------------------

// SetLevel sets the current log level.
func (s *UniLog) SetLevel(level logger_models.Level) {
	s.Logger.SetLevel(level)
}

// -------------------------------------------------------------------------

// Close closes the underlying logger.
func (s *UniLog) Close() {
	if s.Logger != nil {
		s.Logger.Close()
	}
}
