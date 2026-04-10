package logger

import (
	"runtime"

	"github.com/Bastien-Antigravity/universal-logger/src/utils"

	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"

	flex_interfaces "github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	"github.com/Bastien-Antigravity/flexible-logger/src/profiles"
)

var _ interfaces.Logger = (*UniLog)(nil)

// UniLog wraps the flexible-logger library.
type UniLog struct {
	Logger     flex_interfaces.Logger
	NotifQueue <-chan *utils.NotifMessage
}

// -------------------------------------------------------------------------

// NewUniLog initializes a new logger service from an existing logger instance.
// Note: This implementation uses a runtime finalizer to automatically call Close()
// when the logger instance is about to be garbage collected.
func NewUniLog(logger flex_interfaces.Logger) *UniLog {
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
	s.Logger.Log(utils.LevelStream, format, args...)
}

// -------------------------------------------------------------------------

// Logon logs a message at Logon level.
func (s *UniLog) Logon(format string, args ...any) {
	s.Logger.Log(utils.LevelLogon, format, args...)
}

// -------------------------------------------------------------------------

// Logout logs a message at Logout level.
func (s *UniLog) Logout(format string, args ...any) {
	s.Logger.Log(utils.LevelLogout, format, args...)
}

// -------------------------------------------------------------------------

// Trade logs a message at Trade level.
func (s *UniLog) Trade(format string, args ...any) {
	s.Logger.Log(utils.LevelTrade, format, args...)
}

// -------------------------------------------------------------------------

// Schedule logs a message at Schedule level.
func (s *UniLog) Schedule(format string, args ...any) {
	s.Logger.Log(utils.LevelSchedule, format, args...)
}

// -------------------------------------------------------------------------

// Report logs a message at Report level.
func (s *UniLog) Report(format string, args ...any) {
	s.Logger.Log(utils.LevelReport, format, args...)
}

// -------------------------------------------------------------------------

// SetLevel sets the current log level.
func (s *UniLog) SetLevel(level utils.Level) {
	s.Logger.SetLevel(level)
}

// -------------------------------------------------------------------------

// Log logs a message at a specific level.
func (s *UniLog) Log(level utils.Level, format string, args ...any) {
	s.Logger.Log(level, format, args...)
}

// -------------------------------------------------------------------------

// SetLocalNotifQueue sets the notification channel for the local notifier.
// It performs a type assertion to find the appropriate wrapper that supports this.
func (s *UniLog) SetLocalNotifQueue(notifChan chan *utils.NotifMessage) {
	if wrapper, ok := s.Logger.(*profiles.NotifLoggerWrapper); ok {
		wrapper.SetLocalNotifQueue(notifChan)
	}
}

// -------------------------------------------------------------------------

// GetNotifQueue returns the internal notification queue for this logger.
// If the notifier was not enabled during Init, this will return nil.
func (s *UniLog) GetNotifQueue() <-chan *utils.NotifMessage {
	return s.NotifQueue
}

// -------------------------------------------------------------------------

// -------------------------------------------------------------------------

// Unwrap returns the underlying flexible-logger instance.
// This is used by internal utilities for high-performance sink access.
func (s *UniLog) Unwrap() any {
	return s.Logger
}

// Close closes the underlying logger.
func (s *UniLog) Close() {
	if s.Logger != nil {
		s.Logger.Close()
	}
}
