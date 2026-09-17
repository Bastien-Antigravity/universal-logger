package logger_test

// =============================================================================
// ESSENTIAL PROCESS: Unit and concurrency tests for UniLog facade metadata, caller routing, and lifecycle.
//
// DATA FLOW:
//   1. Exercises concurrent metadata mutation under race conditions.
//   2. Validates LogWithCaller propagation and metadata formatting.
//
// KEY PARAMETERS:
//   - TestUniLog_MetadataConcurrency: Validates thread-safe map mutation.
//   - TestUniLog_LogWithCaller: Asserts caller information and metadata preservation.
// =============================================================================

import (
	"fmt"
	"sync"
	"testing"

	flex_interfaces "github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	flex_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
	"github.com/Bastien-Antigravity/universal-logger/src/logger"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// -----------------------------------------------------------------------------

type mockFlexLogger struct {
	mu           sync.Mutex
	lastLevel    flex_models.Level
	lastMsg      string
	lastFile     string
	lastLine     string
	lastFunc     string
	lastMod      string
	logCalled    bool
	callerCalled bool
}

func (m *mockFlexLogger) Debug(format string, args ...any)  {}
func (m *mockFlexLogger) Info(format string, args ...any)   {}
func (m *mockFlexLogger) Warning(format string, args ...any) {}
func (m *mockFlexLogger) Error(format string, args ...any)   {}
func (m *mockFlexLogger) Critical(format string, args ...any) {}
func (m *mockFlexLogger) Stream(format string, args ...any) {}
func (m *mockFlexLogger) Logon(format string, args ...any)  {}
func (m *mockFlexLogger) Logout(format string, args ...any) {}
func (m *mockFlexLogger) Trade(format string, args ...any)  {}
func (m *mockFlexLogger) Schedule(format string, args ...any) {}
func (m *mockFlexLogger) Report(format string, args ...any) {}

func (m *mockFlexLogger) Log(level flex_models.Level, format string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastLevel = level
	if len(args) > 0 {
		m.lastMsg = fmt.Sprintf(format, args...)
	} else {
		m.lastMsg = format
	}
	m.logCalled = true
}

func (m *mockFlexLogger) LogWithCaller(level flex_models.Level, msg, file, line, function, module string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastLevel = level
	m.lastMsg = msg
	m.lastFile = file
	m.lastLine = line
	m.lastFunc = function
	m.lastMod = module
	m.callerCalled = true
}

func (m *mockFlexLogger) SetLevel(level flex_models.Level) {}
func (m *mockFlexLogger) GetLevel() flex_models.Level      { return flex_models.LevelInfo }
func (m *mockFlexLogger) SetCallerSkip(skip int)           {}
func (m *mockFlexLogger) Close()                           {}

var _ flex_interfaces.Logger = (*mockFlexLogger)(nil)

// -----------------------------------------------------------------------------

func TestUniLog_MetadataConcurrency(t *testing.T) {
	mock := &mockFlexLogger{}
	ul := logger.NewUniLog(mock)
	defer ul.Close()

	var wg sync.WaitGroup
	workers := 20
	iterations := 100

	// Concurrent writes via AddMetadata
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ul.AddMetadata(fmt.Sprintf("k_%d", workerID), fmt.Sprintf("v_%d", j))
				ul.GetMetadata()
			}
		}(i)
	}

	// Concurrent writes via SetMetadata
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ul.SetMetadata(map[string]string{
					"env":     "test",
					"version": fmt.Sprintf("1.0.%d", j),
				})
			}
		}(i)
	}

	// Concurrent reads via Log / formatMessageWithMetadata
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ul.Info("concurrent log message %d", j)
			}
		}()
	}

	wg.Wait()
}

// -----------------------------------------------------------------------------

func TestUniLog_LogWithCaller_PreservesDetails(t *testing.T) {
	mock := &mockFlexLogger{}
	ul := logger.NewUniLog(mock)
	defer ul.Close()

	ul.AddMetadata("env", "staging")

	ul.LogWithCaller(utils.LevelWarning, "test warning", "service.py", "105", "execute_task", "my_pkg.service")

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if !mock.callerCalled {
		t.Fatal("Expected LogWithCaller to be called on underlying logger")
	}
	if mock.lastFile != "service.py" {
		t.Errorf("Expected file service.py, got %s", mock.lastFile)
	}
	if mock.lastLine != "105" {
		t.Errorf("Expected line 105, got %s", mock.lastLine)
	}
	if mock.lastFunc != "execute_task" {
		t.Errorf("Expected function execute_task, got %s", mock.lastFunc)
	}
	if mock.lastMod != "my_pkg.service" {
		t.Errorf("Expected module my_pkg.service, got %s", mock.lastMod)
	}
	if mock.lastLevel != interfaces.LevelWarning {
		t.Errorf("Expected level WARNING, got %v", mock.lastLevel)
	}
	expectedMsg := "test warning [meta: env=staging]"
	if mock.lastMsg != expectedMsg {
		t.Errorf("Expected message %q, got %q", expectedMsg, mock.lastMsg)
	}
}

// -----------------------------------------------------------------------------

func TestUniLog_MetadataFormatting(t *testing.T) {
	mock := &mockFlexLogger{}
	ul := logger.NewUniLog(mock)
	defer ul.Close()

	// 1. Without metadata
	ul.Info("hello world")
	if mock.lastMsg != "hello world" {
		t.Errorf("Expected 'hello world', got %q", mock.lastMsg)
	}

	// 2. Single metadata key
	ul.AddMetadata("env", "test")
	ul.Info("message two")
	if mock.lastMsg != "message two [meta: env=test]" {
		t.Errorf("Expected 'message two [meta: env=test]', got %q", mock.lastMsg)
	}

	// 3. Multiple sorted keys
	ul.SetMetadata(map[string]string{
		"zeta":  "last",
		"alpha": "first",
	})
	ul.Info("message three")
	expected := "message three [meta: alpha=first zeta=last]"
	if mock.lastMsg != expected {
		t.Errorf("Expected %q, got %q", expected, mock.lastMsg)
	}
}

// -----------------------------------------------------------------------------

func TestUniLog_LogWithCaller_ModuleFallback(t *testing.T) {
	mock := &mockFlexLogger{}
	ul := logger.NewUniLog(mock)
	defer ul.Close()

	ul.AddMetadata("mod", "fallback.module")
	ul.LogWithCaller(utils.LevelInfo, "test fallback", "main.py", "12", "run", "")

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if mock.lastMod != "fallback.module" {
		t.Errorf("Expected module fallback 'fallback.module', got %q", mock.lastMod)
	}
}

