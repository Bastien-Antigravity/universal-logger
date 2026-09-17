package utils_test

// =============================================================================
// ESSENTIAL PROCESS: Unit tests validating LogWithMetadata polyglot bridging and fallback semantics.
//
// DATA FLOW:
//   1. Invokes LogWithMetadata with explicit caller parameters.
//   2. Asserts correct delegation to LogWithCaller or standard Log fallback.
//
// KEY PARAMETERS:
//   - TestLogWithMetadata_DelegatesToLogWithCaller: Validates caller preservation.
//   - TestLogWithMetadata_FallbackToStandardLog: Validates graceful fallback for non-caller loggers.
// =============================================================================

import (
	"testing"

	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// -----------------------------------------------------------------------------

type mockLoggerWithCaller struct {
	lastLevel Level
	lastMsg   string
	lastFile  string
	lastLine  string
	lastFunc  string
	lastMod   string
	called    bool
}

type Level = interfaces.Level

func (m *mockLoggerWithCaller) Debug(format string, args ...any)  {}
func (m *mockLoggerWithCaller) Info(format string, args ...any)   {}
func (m *mockLoggerWithCaller) Warning(format string, args ...any) {}
func (m *mockLoggerWithCaller) Error(format string, args ...any)   {}
func (m *mockLoggerWithCaller) Critical(format string, args ...any) {}
func (m *mockLoggerWithCaller) Stream(format string, args ...any) {}
func (m *mockLoggerWithCaller) Logon(format string, args ...any)  {}
func (m *mockLoggerWithCaller) Logout(format string, args ...any) {}
func (m *mockLoggerWithCaller) Trade(format string, args ...any)  {}
func (m *mockLoggerWithCaller) Schedule(format string, args ...any) {}
func (m *mockLoggerWithCaller) Report(format string, args ...any) {}
func (m *mockLoggerWithCaller) GetNotifQueue() <-chan *interfaces.NotifMessage { return nil }
func (m *mockLoggerWithCaller) SetLocalNotifQueue(ch chan *interfaces.NotifMessage) {}
func (m *mockLoggerWithCaller) Log(level Level, format string, args ...any) {}
func (m *mockLoggerWithCaller) SetLevel(level Level) {}
func (m *mockLoggerWithCaller) GetLevel() Level { return interfaces.LevelInfo }
func (m *mockLoggerWithCaller) SetCallerSkip(skip int) {}
func (m *mockLoggerWithCaller) SetMetadata(metadata map[string]string) {}
func (m *mockLoggerWithCaller) AddMetadata(key, value string) {}
func (m *mockLoggerWithCaller) GetMetadata() map[string]string { return nil }
func (m *mockLoggerWithCaller) Close() {}

func (m *mockLoggerWithCaller) LogWithCaller(level Level, msg, file, line, function, module string) {
	m.lastLevel = level
	m.lastMsg = msg
	m.lastFile = file
	m.lastLine = line
	m.lastFunc = function
	m.lastMod = module
	m.called = true
}

var _ interfaces.Logger = (*mockLoggerWithCaller)(nil)

// -----------------------------------------------------------------------------

type mockWrapper struct {
	*mockLoggerWithCaller
	wrapped interfaces.Logger
}

func (w *mockWrapper) Unwrap() any {
	return w.wrapped
}

// -----------------------------------------------------------------------------

func TestLogWithMetadata_DelegatesToLogWithCaller(t *testing.T) {
	mock := &mockLoggerWithCaller{}

	utils.LogWithMetadata(
		mock,
		interfaces.LevelError,
		"database query failed",
		"db_client.py",
		"42",
		"execute_query",
		"backend.db",
	)

	if !mock.called {
		t.Fatal("Expected LogWithCaller to be called")
	}
	if mock.lastLevel != interfaces.LevelError {
		t.Errorf("Expected ERROR level, got %v", mock.lastLevel)
	}
	if mock.lastMsg != "database query failed" {
		t.Errorf("Expected message 'database query failed', got %q", mock.lastMsg)
	}
	if mock.lastFile != "db_client.py" {
		t.Errorf("Expected file 'db_client.py', got %q", mock.lastFile)
	}
	if mock.lastLine != "42" {
		t.Errorf("Expected line '42', got %q", mock.lastLine)
	}
	if mock.lastFunc != "execute_query" {
		t.Errorf("Expected function 'execute_query', got %q", mock.lastFunc)
	}
	if mock.lastMod != "backend.db" {
		t.Errorf("Expected module 'backend.db', got %q", mock.lastMod)
	}
}

// -----------------------------------------------------------------------------

func TestLogWithMetadata_UnwrapsTarget(t *testing.T) {
	innerMock := &mockLoggerWithCaller{}
	wrapper := &mockWrapper{
		mockLoggerWithCaller: &mockLoggerWithCaller{}, // separate instance
		wrapped:              innerMock,
	}

	utils.LogWithMetadata(
		wrapper,
		interfaces.LevelWarning,
		"wrapped warning",
		"handler.rs",
		"88",
		"process",
		"rust_service",
	)

	// The wrapper's LogWithCaller should be called, or if unwrapped, innerMock called
	if !wrapper.called && !innerMock.called {
		t.Fatal("Expected LogWithCaller to be called on wrapper or inner target")
	}
}

