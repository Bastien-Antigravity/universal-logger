package interfaces

// =============================================================================
// ESSENTIAL PROCESS: Core Logger facade interface definition establishing standard logging methods and caller injection.
//
// DATA FLOW:
//   1. Defines contract for 12 log severity methods.
//   2. Declares LogWithCaller for explicit caller stack metadata propagation.
//   3. Specifies metadata inspection, level adjustment, and lifecycle closure contracts.
//
// KEY PARAMETERS:
//   - Logger: Unified logging contract for the Bastien-Antigravity ecosystem.
// =============================================================================


// Logger is the main interface for logging across the Bastien-Antigravity ecosystem.
// It is a facade that ensures microservices remain decoupled from the underlying logging engine.
type Logger interface {
	// -------------------------------------------------------------------------
	// Core Logging methods (Printf-style)
	// -------------------------------------------------------------------------

	// Debug logs a message at Debug level.
	Debug(format string, args ...any)

	// Info logs a message at Info level.
	Info(format string, args ...any)

	// Warning logs a message at Warning level.
	Warning(format string, args ...any)

	// Error logs a message at Error level.
	Error(format string, args ...any)

	// Critical logs a message at Critical level.
	Critical(format string, args ...any)

	// -------------------------------------------------------------------------
	// Specialized Domain Logging
	// -------------------------------------------------------------------------

	// Stream logs a message at Stream level.
	Stream(format string, args ...any)

	// Logon logs a message at Logon level.
	Logon(format string, args ...any)

	// Logout logs a message at Logout level.
	Logout(format string, args ...any)

	// Trade logs a message at Trade level.
	Trade(format string, args ...any)

	// Schedule logs a message at Schedule level.
	Schedule(format string, args ...any)

	// Report logs a message at Report level.
	Report(format string, args ...any)

	// -------------------------------------------------------------------------
	// Generic and Control Methods
	// -------------------------------------------------------------------------

	// GetNotifQueue returns the internal notification queue for this logger.
	// If the notifier was not enabled during Init, this will return nil.
	GetNotifQueue() <-chan *NotifMessage

	// SetLocalNotifQueue allows manual binding of a notification channel.
	// This is useful for services like notif-server that manage their own channels.
	SetLocalNotifQueue(notifChan chan *NotifMessage)

	// Log logs a message at a specific level.
	Log(level Level, format string, args ...any)

	// LogWithCaller logs a message with explicit caller stack metadata.
	LogWithCaller(level Level, msg, file, line, function, module string)

	// SetLevel sets the current log level.
	SetLevel(level Level)

	// GetLevel returns the current log level.
	GetLevel() Level

	// SetCallerSkip sets the number of stack frames to skip when detecting source info.
	SetCallerSkip(skip int)

	// -------------------------------------------------------------------------
	// Metadata and Tagging
	// -------------------------------------------------------------------------

	// SetMetadata replaces all existing metadata with the provided map.
	SetMetadata(metadata map[string]string)

	// AddMetadata adds a single key-value pair to the logger's metadata.
	AddMetadata(key, value string)

	// GetMetadata returns a copy of the current logger metadata.
	GetMetadata() map[string]string

	// Close flushes any buffered logs and closes the handler.
	Close()
}
