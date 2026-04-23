package bootstrap

import (
	"testing"
	"time"

	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

func TestLogLevelSynchronization(t *testing.T) {
	// 1. Initialize with INFO level
	distConfig, uniLog := Init("sync-test", "standalone", "devel", "INFO", false, nil)
	defer uniLog.Close()

	if uniLog.GetLevel() != utils.LevelInfo {
		t.Errorf("Expected initial level INFO, got %v", uniLog.GetLevel())
	}

	// 2. Simulate configuration update (logger.level = DEBUG)
	distConfig.SetConfig("logger", "level", "DEBUG")

	// Give a small amount of time for the callback to propagate
	time.Sleep(100 * time.Millisecond)

	if uniLog.GetLevel() != utils.LevelDebug {
		t.Errorf("Expected level to sync to DEBUG, got %v", uniLog.GetLevel())
	}
}

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
