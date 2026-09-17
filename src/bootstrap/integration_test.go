package bootstrap

// =============================================================================
// ESSENTIAL PROCESS: Integration tests validating configuration injection, dynamic log level synchronization, and notification lifecycle.
//
// DATA FLOW:
//   1. Injects configuration updates and asserts live log level synchronization.
//   2. Attaches local notification queues and verifies severity threshold gating.
//   3. Stresses notification buffer to guarantee non-blocking delivery on saturation.
//
// KEY PARAMETERS:
//   - TestLogLevelSynchronization: Validates basic config-driven level update.
//   - TestDynamicLogLevelFiltering_ConfigurationSync: Validates dynamic multi-level transitions.
//   - TestLocalNotification_LevelThresholdAndPayload: Asserts alert gating across all log levels.
//   - TestLocalNotification_NonBlockingOnSaturation: Verifies non-blocking buffer overflow behavior.
// =============================================================================


import (
	"testing"
	"time"

	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// -----------------------------------------------------------------------------

func TestLogLevelSynchronization(t *testing.T) {
	// 1. Initialize with INFO level
	distConfig, uniLog := Init("sync-test", "standalone", "devel", "INFO", false, nil)
	defer uniLog.Close()

	if uniLog.GetLevel() != utils.LevelInfo {
		t.Errorf("Expected initial level INFO, got %v", uniLog.GetLevel())
	}

	// 2. Simulate configuration update (logger.level = DEBUG)
	if err := distConfig.SetConfig("logger", "level", "DEBUG"); err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Give a small amount of time for the callback to propagate
	time.Sleep(100 * time.Millisecond)

	if uniLog.GetLevel() != utils.LevelDebug {
		t.Errorf("Expected level to sync to DEBUG, got %v", uniLog.GetLevel())
	}
}

// -----------------------------------------------------------------------------

func TestDynamicLogLevelFiltering_ConfigurationSync(t *testing.T) {
	distConfig, uniLog := Init("dyn-filter-sync", "standalone", "devel", "INFO", false, nil)
	defer uniLog.Close()

	if uniLog.GetLevel() != utils.LevelInfo {
		t.Fatalf("Expected initial level INFO, got %v", uniLog.GetLevel())
	}

	// 1. Update config to ERROR
	if err := distConfig.SetConfig("logger", "level", "ERROR"); err != nil {
		t.Fatalf("Failed to set config to ERROR: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	if uniLog.GetLevel() != utils.LevelError {
		t.Errorf("Expected level to sync to ERROR, got %v", uniLog.GetLevel())
	}

	// 2. Update config to WARNING
	if err := distConfig.SetConfig("logger", "level", "WARNING"); err != nil {
		t.Fatalf("Failed to set config to WARNING: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	if uniLog.GetLevel() != utils.LevelWarning {
		t.Errorf("Expected level to sync to WARNING, got %v", uniLog.GetLevel())
	}

	// 3. Update config to DEBUG
	if err := distConfig.SetConfig("logger", "level", "DEBUG"); err != nil {
		t.Fatalf("Failed to set config to DEBUG: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	if uniLog.GetLevel() != utils.LevelDebug {
		t.Errorf("Expected level to sync to DEBUG, got %v", uniLog.GetLevel())
	}
}

// -----------------------------------------------------------------------------

func TestMetadataInjection(t *testing.T) {
	meta := map[string]string{
		"env":     "test",
		"version": "1.1.7",
	}

	opts := BootstrapOptions{
		Name:            "meta-test",
		ConfigProfile:   "standalone",
		LoggerProfile:   "devel",
		InitialLogLevel: utils.LevelInfo,
		Metadata:        meta,
	}

	_, uniLog := InitWithOptions(opts)
	defer uniLog.Close()

	if uniLog == nil {
		t.Fatal("Expected logger to be initialized")
	}
}

// -----------------------------------------------------------------------------

func TestConfigInjection(t *testing.T) {
	// 1. Create a config instance manually
	baseConfig, _ := Init("base", "standalone", "devel", "INFO", false, nil)

	// 2. Inject it into a second session
	opts := BootstrapOptions{
		Name:           "injected",
		LoggerProfile:  "devel",
		ExistingConfig: baseConfig,
	}

	distConfig, _ := InitWithOptions(opts)

	if distConfig != baseConfig {
		t.Fatal("Expected injected config to be used, but got a new instance")
	}
}

// -----------------------------------------------------------------------------

func TestManualNotifierBinding(t *testing.T) {
	// 1. Initialize with useLocalNotifier = true on DEVEL profile (which supports it safely)
	_, uniLog := Init("notif-test", "standalone", "devel", "INFO", true, nil)
	defer uniLog.Close()

	// 2. Create a manual channel
	myChan := make(chan *utils.NotifMessage, 10)
	uniLog.SetLocalNotifQueue(myChan)

	// 3. Log a warning (which triggers notification in most profiles if notifier is set)
	uniLog.Warning("Test notification")

	// 4. Check if it arrived in our channel
	select {
	case msg := <-myChan:
		if msg.Message != "Test notification" {
			t.Errorf("Expected 'Test notification', got %s", msg.Message)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for notification")
	}
}

// -----------------------------------------------------------------------------

func TestLocalNotification_LevelThresholdAndPayload(t *testing.T) {
	// 1. Initialize with useLocalNotifier = true
	_, uniLog := Init("thresh-test", "standalone", "devel", "DEBUG", true, nil)
	defer uniLog.Close()

	notifQueue := uniLog.GetNotifQueue()
	if notifQueue == nil {
		t.Fatal("Expected non-nil NotifQueue from GetNotifQueue()")
	}

	// 2. Debug and Info must NOT emit alerts to the notification queue
	uniLog.Debug("suppressed debug alert")
	uniLog.Info("suppressed info alert")
	uniLog.Stream("suppressed stream alert")
	uniLog.LogWithCaller(utils.LevelInfo, "suppressed caller alert", "main.go", "10", "run", "core")
	utils.LogWithMetadata(uniLog, utils.LevelInfo, "suppressed meta alert", "main.go", "20", "run", "core")

	select {
	case msg := <-notifQueue:
		t.Fatalf("Expected zero notifications for low-severity logs, got: %+v", msg)
	case <-time.After(100 * time.Millisecond):
		// Expected: queue remains empty
	}

	// 3. Warning must trigger an alert with Level == "WARNING"
	uniLog.Warning("warning alert message")
	select {
	case msg := <-notifQueue:
		if msg.Level != "WARNING" {
			t.Errorf("Expected level WARNING, got %q", msg.Level)
		}
		if msg.Message != "warning alert message" {
			t.Errorf("Expected 'warning alert message', got %q", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for Warning alert")
	}

	// 4. Error must trigger an alert with Level == "ERROR"
	uniLog.Error("error alert message")
	select {
	case msg := <-notifQueue:
		if msg.Level != "ERROR" {
			t.Errorf("Expected level ERROR, got %q", msg.Level)
		}
		if msg.Message != "error alert message" {
			t.Errorf("Expected 'error alert message', got %q", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for Error alert")
	}

	// 5. Critical must trigger an alert with Level == "CRITICAL"
	uniLog.Critical("critical alert message")
	select {
	case msg := <-notifQueue:
		if msg.Level != "CRITICAL" {
			t.Errorf("Expected level CRITICAL, got %q", msg.Level)
		}
		if msg.Message != "critical alert message" {
			t.Errorf("Expected 'critical alert message', got %q", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for Critical alert")
	}

	// 6. LogWithCaller and LogWithMetadata high-severity alerts
	uniLog.LogWithCaller(utils.LevelWarning, "caller warning alert", "service.py", "55", "process", "worker")
	select {
	case msg := <-notifQueue:
		if msg.Level != "WARNING" {
			t.Errorf("Expected level WARNING, got %q", msg.Level)
		}
		if msg.Message != "caller warning alert" {
			t.Errorf("Expected 'caller warning alert', got %q", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for LogWithCaller alert")
	}

	utils.LogWithMetadata(uniLog, utils.LevelError, "meta error alert", "app.cpp", "120", "on_event", "engine")
	select {
	case msg := <-notifQueue:
		if msg.Level != "ERROR" {
			t.Errorf("Expected level ERROR, got %q", msg.Level)
		}
		if msg.Message != "meta error alert" {
			t.Errorf("Expected 'meta error alert', got %q", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for LogWithMetadata alert")
	}
}

// -----------------------------------------------------------------------------

func TestLocalNotification_NonBlockingOnSaturation(t *testing.T) {
	_, uniLog := Init("saturation-test", "standalone", "devel", "INFO", true, nil)
	defer uniLog.Close()

	notifQueue := uniLog.GetNotifQueue()
	if notifQueue == nil {
		t.Fatal("Expected non-nil NotifQueue")
	}

	// Queue buffer capacity is 1024.
	// Emit 1200 warning logs without draining the queue.
	// This MUST NOT block, deadlock, or panic.
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1200; i++ {
			uniLog.Warning("flood-warning-%d", i)
		}
		close(done)
	}()

	select {
	case <-done:
		// Succeeded within timeout
	case <-time.After(2 * time.Second):
		t.Fatal("Logging blocked on saturated notification queue!")
	}

	// Drain the queue: exactly 1024 messages should be readable
	drained := 0
	for {
		select {
		case msg := <-notifQueue:
			if msg == nil {
				t.Fatal("Received nil message from notifQueue")
			}
			drained++
		default:
			goto DRAINED
		}
	}
DRAINED:
	if drained != 1024 {
		t.Errorf("Expected exactly 1024 messages in queue before drop, got %d", drained)
	}

	// After draining, new notifications should be received normally again
	uniLog.Warning("recovered-warning")
	select {
	case msg := <-notifQueue:
		if msg.Message != "recovered-warning" {
			t.Errorf("Expected 'recovered-warning', got %s", msg.Message)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for notification after queue was drained")
	}
}

