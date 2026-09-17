package logger_test

// =============================================================================
// ESSENTIAL PROCESS: Unit and concurrency tests for UniLog facade metadata, caller routing, dynamic level filtering, and notification thresholds.
//
// DATA FLOW:
//   1. Exercises concurrent metadata mutation under race conditions.
//   2. Validates LogWithCaller propagation and metadata formatting.
//   3. Tests dynamic log level filtering across standard and caller-injected paths.
//   4. Tests alert notification dispatching and severity threshold gating.
//
// KEY PARAMETERS:
//   - TestUniLog_MetadataConcurrency: Validates thread-safe map mutation.
//   - TestUniLog_LogWithCaller: Asserts caller information and metadata preservation.
//   - TestUniLog_DynamicLevelFiltering_Unit: Validates dynamic sink filtering on level adjustments.
//   - TestUniLog_NotificationThresholds_Unit: Asserts alerts fire only on Warning/Error/Critical.
// =============================================================================

import (
	"fmt"
	"sync"
	"testing"
	"time"

	flex_engine "github.com/Bastien-Antigravity/flexible-logger/src/engine"
	flex_factory "github.com/Bastien-Antigravity/flexible-logger/src/factory"
	flex_interfaces "github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
	flex_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
	flex_notifier "github.com/Bastien-Antigravity/flexible-logger/src/notifier"
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

// -----------------------------------------------------------------------------

type capturedRecord struct {
	Level   flex_models.Level
	Message string
	File    string
	Line    string
	Func    string
	Module  string
}

type capturingSink struct {
	mu      sync.Mutex
	entries []capturedRecord
}

func (s *capturingSink) Write(entry *flex_models.LogEntry) error {
	defer entry.Release()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, capturedRecord{
		Level:   entry.Level,
		Message: entry.Message,
		File:    entry.Filename,
		Line:    entry.LineNumber,
		Func:    entry.FunctionName,
		Module:  entry.Module,
	})
	return nil
}

func (s *capturingSink) Close() error {
	return nil
}

func (s *capturingSink) getEntries() []capturedRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := make([]capturedRecord, len(s.entries))
	copy(copied, s.entries)
	return copied
}

func (s *capturingSink) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = nil
}

// -----------------------------------------------------------------------------

func TestUniLog_DynamicLevelFiltering_Unit(t *testing.T) {
	sink := &capturingSink{}
	flexEngine := flex_factory.CreateLogEngine("unit-filter-test", flex_models.LevelInfo, sink, false, 0).(*flex_engine.LogEngine)
	ul := logger.NewUniLog(flexEngine)
	defer ul.Close()

	// 1. Initial level is INFO
	if ul.GetLevel() != utils.LevelInfo {
		t.Fatalf("Expected initial level INFO, got %v", ul.GetLevel())
	}

	ul.Debug("debug-1")
	ul.Info("info-1")
	ul.Warning("warning-1")
	ul.Error("error-1")
	ul.Critical("critical-1")

	entries := sink.getEntries()
	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries (Debug dropped), got %d: %+v", len(entries), entries)
	}
	if entries[0].Message != "info-1" || entries[1].Message != "warning-1" || entries[2].Message != "error-1" || entries[3].Message != "critical-1" {
		t.Errorf("Unexpected entries at INFO level: %+v", entries)
	}

	// 2. Dynamically elevate level to WARNING
	ul.SetLevel(utils.LevelWarning)
	if ul.GetLevel() != utils.LevelWarning {
		t.Fatalf("Expected level to be updated to WARNING, got %v", ul.GetLevel())
	}

	sink.clear()
	ul.Debug("debug-2")
	ul.Info("info-2")
	ul.Warning("warning-2")
	ul.Error("error-2")

	// Also test LogWithCaller and LogWithMetadata at WARNING threshold
	ul.LogWithCaller(utils.LevelInfo, "caller-info-filtered", "test.py", "10", "run", "pkg")
	ul.LogWithCaller(utils.LevelWarning, "caller-warning-allowed", "test.py", "20", "run", "pkg")
	utils.LogWithMetadata(ul, utils.LevelInfo, "meta-info-filtered", "main.go", "30", "exec", "core")
	utils.LogWithMetadata(ul, utils.LevelError, "meta-error-allowed", "main.go", "40", "exec", "core")

	entries = sink.getEntries()
	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries at WARNING level (Debug/Info dropped), got %d: %+v", len(entries), entries)
	}
	expectedMsgs := []string{"warning-2", "error-2", "caller-warning-allowed", "meta-error-allowed"}
	for i, exp := range expectedMsgs {
		if entries[i].Message != exp {
			t.Errorf("Entry %d expected message %q, got %q", i, exp, entries[i].Message)
		}
	}

	// 3. Dynamically lower level to DEBUG
	ul.SetLevel(utils.LevelDebug)
	if ul.GetLevel() != utils.LevelDebug {
		t.Fatalf("Expected level to be updated to DEBUG, got %v", ul.GetLevel())
	}

	sink.clear()
	ul.Debug("debug-3")
	ul.LogWithCaller(utils.LevelDebug, "caller-debug-3", "test.py", "50", "run", "pkg")
	utils.LogWithMetadata(ul, utils.LevelDebug, "meta-debug-3", "main.go", "60", "exec", "core")

	entries = sink.getEntries()
	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries at DEBUG level, got %d: %+v", len(entries), entries)
	}
	if entries[0].Message != "debug-3" || entries[1].Message != "caller-debug-3" || entries[2].Message != "meta-debug-3" {
		t.Errorf("Unexpected entries at DEBUG level: %+v", entries)
	}
}

// -----------------------------------------------------------------------------

func TestUniLog_NotificationThresholds_Unit(t *testing.T) {
	sink := &capturingSink{}
	flexEngine := flex_factory.CreateLogEngine("unit-notif-test", flex_models.LevelDebug, sink, false, 0).(*flex_engine.LogEngine)
	localNotif := flex_notifier.NewLocalNotifier()
	notifChan := make(chan *flex_models.NotifMessage, 100)
	localNotif.SetQueue(notifChan)
	flexEngine.Notifier = localNotif

	ul := logger.NewUniLog(flexEngine)
	defer ul.Close()

	// 1. Low-severity and domain logs should NOT trigger notifications
	ul.Debug("dbg-msg")
	ul.Info("info-msg")
	ul.Stream("stream-msg")
	ul.Logon("logon-msg")
	ul.Logout("logout-msg")
	ul.Trade("trade-msg")
	ul.Schedule("schedule-msg")
	ul.Report("report-msg")
	ul.LogWithCaller(utils.LevelInfo, "caller-info", "a.py", "1", "f", "m")
	utils.LogWithMetadata(ul, utils.LevelInfo, "meta-info", "b.go", "2", "f", "m")

	select {
	case msg := <-notifChan:
		t.Fatalf("Expected no notification for low-severity logs, got: %+v", msg)
	case <-time.After(50 * time.Millisecond):
		// Expected: no notification received
	}

	// 2. High-severity logs (Warning, Error, Critical) MUST trigger notifications
	ul.Warning("warn-msg")
	select {
	case msg := <-notifChan:
		if msg.Level != "WARNING" || msg.Message != "warn-msg" {
			t.Errorf("Unexpected warning notification: %+v", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for warning notification")
	}

	ul.Error("error-msg")
	select {
	case msg := <-notifChan:
		if msg.Level != "ERROR" || msg.Message != "error-msg" {
			t.Errorf("Unexpected error notification: %+v", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for error notification")
	}

	ul.Critical("crit-msg")
	select {
	case msg := <-notifChan:
		if msg.Level != "CRITICAL" || msg.Message != "crit-msg" {
			t.Errorf("Unexpected critical notification: %+v", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for critical notification")
	}

	// 3. LogWithCaller and LogWithMetadata high-severity alerts
	ul.LogWithCaller(utils.LevelWarning, "caller-warn-msg", "svc.go", "42", "run", "svc")
	select {
	case msg := <-notifChan:
		if msg.Level != "WARNING" || msg.Message != "caller-warn-msg" {
			t.Errorf("Unexpected LogWithCaller notification: %+v", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for LogWithCaller notification")
	}

	utils.LogWithMetadata(ul, utils.LevelCritical, "meta-crit-msg", "svc.py", "88", "handler", "api")
	select {
	case msg := <-notifChan:
		if msg.Level != "CRITICAL" || msg.Message != "meta-crit-msg" {
			t.Errorf("Unexpected LogWithMetadata notification: %+v", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timed out waiting for LogWithMetadata notification")
	}
}


