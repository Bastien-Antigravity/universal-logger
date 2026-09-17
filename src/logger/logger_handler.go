package logger

// =============================================================================
// ESSENTIAL PROCESS: Unified logger facade coordinating flexible-logger delegates, thread-safe metadata enrichment, and polyglot caller dispatch.
//
// DATA FLOW:
//   1. Receives logging calls across standard, domain, or polyglot FFI paths.
//   2. Enriches payloads with thread-safe metadata key-value tags.
//   3. Forwards structured entries to underlying flexible-logger engine.
//
// KEY PARAMETERS:
//   - UniLog: Primary facade wrapper around flexible-logger Logger.
//   - metaMu: RWMutex guarding concurrent metadata read/write operations.
// =============================================================================

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"

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
	Metadata   map[string]string
	metaMu     sync.RWMutex
	closeOnce  sync.Once
}

// -----------------------------------------------------------------------------

// NewUniLog initializes a new logger service from an existing logger instance.
// Note: This implementation uses a runtime finalizer to automatically call Close()
// when the logger instance is about to be garbage collected.
func NewUniLog(logger flex_interfaces.Logger) *UniLog {
	res := &UniLog{
		Logger:   logger,
		Metadata: make(map[string]string),
	}

	// Initialize with 1 to skip this facade layer by default
	res.Logger.SetCallerSkip(1)

	// Register finalizer for automatic cleanup
	runtime.SetFinalizer(res, func(ul *UniLog) {
		ul.Close()
	})

	return res
}

// -----------------------------------------------------------------------------

func (s *UniLog) formatMessageWithMetadata(msg string) string {
	s.metaMu.RLock()
	defer s.metaMu.RUnlock()
	if len(s.Metadata) == 0 {
		return msg
	}
	keys := make([]string, 0, len(s.Metadata))
	for k := range s.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tags := make([]string, 0, len(keys))
	for _, k := range keys {
		tags = append(tags, fmt.Sprintf("%s=%s", k, s.Metadata[k]))
	}
	return fmt.Sprintf("%s [meta: %s]", msg, strings.Join(tags, " "))
}

// -----------------------------------------------------------------------------

// SetMetadata replaces all existing metadata with the provided map.
func (s *UniLog) SetMetadata(metadata map[string]string) {
	s.metaMu.Lock()
	defer s.metaMu.Unlock()
	s.Metadata = make(map[string]string, len(metadata))
	for k, v := range metadata {
		s.Metadata[k] = v
	}
}

// -----------------------------------------------------------------------------

// AddMetadata adds a single key-value pair to the logger's metadata.
func (s *UniLog) AddMetadata(key, value string) {
	s.metaMu.Lock()
	defer s.metaMu.Unlock()
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
}

// -----------------------------------------------------------------------------

// GetMetadata returns a thread-safe copy of the logger's metadata.
func (s *UniLog) GetMetadata() map[string]string {
	s.metaMu.RLock()
	defer s.metaMu.RUnlock()
	if s.Metadata == nil {
		return nil
	}
	res := make(map[string]string, len(s.Metadata))
	for k, v := range s.Metadata {
		res[k] = v
	}
	return res
}

// -----------------------------------------------------------------------------

// Debug logs a message at Debug level.
func (s *UniLog) Debug(format string, args ...any) {
	s.Log(utils.LevelDebug, format, args...)
}

// -----------------------------------------------------------------------------

// Info logs a message at Info level.
func (s *UniLog) Info(format string, args ...any) {
	s.Log(utils.LevelInfo, format, args...)
}

// -----------------------------------------------------------------------------

// Warning logs a message at Warning level.
func (s *UniLog) Warning(format string, args ...any) {
	s.Log(utils.LevelWarning, format, args...)
}

// -----------------------------------------------------------------------------

// Error logs a message at Error level.
func (s *UniLog) Error(format string, args ...any) {
	s.Log(utils.LevelError, format, args...)
}

// -----------------------------------------------------------------------------

// Critical logs a message at Critical level.
func (s *UniLog) Critical(format string, args ...any) {
	s.Log(utils.LevelCritical, format, args...)
}

// -----------------------------------------------------------------------------

// Stream logs a message at Stream level.
func (s *UniLog) Stream(format string, args ...any) {
	s.Log(utils.LevelStream, format, args...)
}

// -----------------------------------------------------------------------------

// Logon logs a message at Logon level.
func (s *UniLog) Logon(format string, args ...any) {
	s.Log(utils.LevelLogon, format, args...)
}

// -----------------------------------------------------------------------------

// Logout logs a message at Logout level.
func (s *UniLog) Logout(format string, args ...any) {
	s.Log(utils.LevelLogout, format, args...)
}

// -----------------------------------------------------------------------------

// Trade logs a message at Trade level.
func (s *UniLog) Trade(format string, args ...any) {
	s.Log(utils.LevelTrade, format, args...)
}

// -----------------------------------------------------------------------------

// Schedule logs a message at Schedule level.
func (s *UniLog) Schedule(format string, args ...any) {
	s.Log(utils.LevelSchedule, format, args...)
}

// -----------------------------------------------------------------------------

// Report logs a message at Report level.
func (s *UniLog) Report(format string, args ...any) {
	s.Log(utils.LevelReport, format, args...)
}

// -----------------------------------------------------------------------------

// SetLevel sets the current log level.
func (s *UniLog) SetLevel(level utils.Level) {
	s.Logger.SetLevel(level)
}

// -----------------------------------------------------------------------------

// GetLevel returns the current log level.
func (s *UniLog) GetLevel() utils.Level {
	return s.Logger.GetLevel()
}

// -----------------------------------------------------------------------------

// Log logs a message at a specific level.
func (s *UniLog) Log(level utils.Level, format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	msg = s.formatMessageWithMetadata(msg)
	s.Logger.Log(level, "%s", msg)
}

// -----------------------------------------------------------------------------

// LogWithCaller logs a message with explicit caller stack metadata.
func (s *UniLog) LogWithCaller(level utils.Level, msg, file, line, function, module string) {
	msg = s.formatMessageWithMetadata(msg)
	if module == "" {
		s.metaMu.RLock()
		if m, ok := s.Metadata["mod"]; ok {
			module = m
		} else if m, ok := s.Metadata["module"]; ok {
			module = m
		}
		s.metaMu.RUnlock()
	}
	s.Logger.LogWithCaller(level, msg, file, line, function, module)
}

// -----------------------------------------------------------------------------

// SetCallerSkip sets the number of stack frames to skip when detecting source info.
// It automatically adds 1 to the provided skip value to account for this facade layer.
func (s *UniLog) SetCallerSkip(skip int) {
	s.Logger.SetCallerSkip(skip + 1)
}

// -----------------------------------------------------------------------------

// SetLocalNotifQueue sets the notification channel for the local notifier.
// It performs a type assertion to find the appropriate wrapper that supports this.
func (s *UniLog) SetLocalNotifQueue(notifChan chan *interfaces.NotifMessage) {
	if wrapper, ok := s.Logger.(*profiles.NotifLoggerWrapper); ok {
		wrapper.SetLocalNotifQueue(notifChan)
	}
}

// -----------------------------------------------------------------------------

// GetNotifQueue returns the internal notification queue for this logger.
// If the notifier was not enabled during Init, this will return nil.
func (s *UniLog) GetNotifQueue() <-chan *utils.NotifMessage {
	return s.NotifQueue
}

// -----------------------------------------------------------------------------

// Unwrap returns the underlying flexible-logger instance.
// This is used by internal utilities for high-performance sink access.
func (s *UniLog) Unwrap() any {
	return s.Logger
}

// -----------------------------------------------------------------------------

// Close closes the underlying logger idempotently and unregisters the GC finalizer.
func (s *UniLog) Close() {
	s.closeOnce.Do(func() {
		runtime.SetFinalizer(s, nil)
		if s.Logger != nil {
			s.Logger.Close()
		}
	})
}
