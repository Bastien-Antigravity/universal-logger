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
	"testing"
)

func TestInitWithLocalNotifier(t *testing.T) {
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
