package bootstrap

// =============================================================================
// ESSENTIAL PROCESS: Unit tests validating bootstrap initialization and default profile behaviors.
//
// DATA FLOW:
//   1. Executes bootstrap routines under test configurations.
//   2. Asserts non-nil sessions, default log levels, and clean shutdown.
//
// KEY PARAMETERS:
//   - TestInitService: Validates zero-config bootstrap entrypoint.
// =============================================================================

import (
	"os"
	"path/filepath"
	"testing"
)

func setupIsolatedBootstrapConfig(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("LOG_DIR", filepath.Join(tmpDir, "logs"))
	cfgPath := filepath.Join(tmpDir, "standalone.yaml")
	yamlContent := `common:
  name: "unilog-test"
  public_ip: "127.0.0.1"
`
	if err := os.WriteFile(cfgPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatalf("Failed to write temp standalone config: %v", err)
	}
	t.Setenv("CONFIG_PATH", cfgPath)
}

func TestInitWithLocalNotifier(t *testing.T) {
	setupIsolatedBootstrapConfig(t)
	// Initialize with useLocalNotifier = true
	_, uniLog := Init("test-app", "standalone", "devel", "INFO", true, nil)
	defer uniLog.Close()

	if uniLog.GetNotifQueue() == nil {
		t.Fatal("Expected NotifQueue to be initialized when useLocalNotifier is true")
	}

	queue := uniLog.GetNotifQueue()
	if queue == nil {
		t.Fatal("Expected GetNotifQueue to return the initialized channel")
	}

	// Verify buffer size
	if cap(queue) != 1024 {
		t.Errorf("Expected NotifQueue buffer size to be 1024, got %d", cap(queue))
	}
}

func TestInitWithoutLocalNotifier(t *testing.T) {
	setupIsolatedBootstrapConfig(t)
	// Initialize with useLocalNotifier = false
	_, uniLog := Init("test-app", "standalone", "devel", "INFO", false, nil)
	defer uniLog.Close()

	if uniLog.GetNotifQueue() != nil {
		t.Fatal("Expected NotifQueue to be nil when useLocalNotifier is false")
	}

	if uniLog.GetNotifQueue() != nil {
		t.Fatal("Expected GetNotifQueue to return nil when not enabled")
	}
}
