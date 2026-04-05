package bootstrap

import (
	"testing"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

func TestInitWithLocalNotifier(t *testing.T) {
	// Initialize with useLocalNotifier = true
	_, uniLog := Init("test-app", "standalone", "devel", utils.LevelInfo, true)
	defer uniLog.Close()

	if uniLog.NotifQueue == nil {
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
	_, uniLog := Init("test-app", "standalone", "devel", utils.LevelInfo, false)
	defer uniLog.Close()

	if uniLog.NotifQueue != nil {
		t.Fatal("Expected NotifQueue to be nil when useLocalNotifier is false")
	}

	if uniLog.GetNotifQueue() != nil {
		t.Fatal("Expected GetNotifQueue to return nil when not enabled")
	}
}
