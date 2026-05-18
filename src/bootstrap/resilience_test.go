package bootstrap

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Bastien-Antigravity/safe-socket/src/facade"
	"github.com/Bastien-Antigravity/safe-socket/src/interfaces"
	"github.com/Bastien-Antigravity/safe-socket/src/models"
	"github.com/Bastien-Antigravity/safe-socket/src/profiles"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// MockServer simulates a remote service using the safe-socket library.
type MockServer struct {
	server   *facade.SocketServer
	received chan string
	shutdown chan struct{}
	port     string
	mu       sync.Mutex
	running  bool
	conns    []interfaces.TransportConnection
}

func NewMockServer(t *testing.T, name string, profileType string) (*MockServer, string, string) {
	t.Helper()

	var profile interfaces.SocketProfile
	if profileType == "tcp-hello" {
		profile = profiles.NewTcpHelloServerProfile(name, "127.0.0.1:0", 5000)
	} else {
		profile = profiles.NewTcpServerProfile(name, "127.0.0.1:0", 5000)
	}

	config := models.SocketConfig{
		Deadline: 30 * time.Second,
	}

	server := facade.NewSocketServer(profile, config)
	if err := server.Listen(); err != nil {
		t.Fatalf("Failed to start safe-socket MockServer [%s]: %v", name, err)
	}

	addr, err := server.GetAddr()
	if err != nil {
		t.Fatalf("Failed to get address for MockServer [%s]: %v", name, err)
	}
	host, port, _ := net.SplitHostPort(addr)

	m := &MockServer{
		server:   server,
		received: make(chan string, 100),
		shutdown: make(chan struct{}),
		port:     port,
		running:  true,
	}

	go m.acceptLoop()
	return m, host, port
}

// NewMockServerOnPort creates a MockServer bound to a specific port.
func NewMockServerOnPort(t *testing.T, name string, profileType string, port string) *MockServer {
	t.Helper()

	var profile interfaces.SocketProfile
	addr := "127.0.0.1:" + port
	if profileType == "tcp-hello" {
		profile = profiles.NewTcpHelloServerProfile(name, addr, 5000)
	} else {
		profile = profiles.NewTcpServerProfile(name, addr, 5000)
	}

	config := models.SocketConfig{
		Deadline: 30 * time.Second,
	}

	server := facade.NewSocketServer(profile, config)
	if err := server.Listen(); err != nil {
		t.Fatalf("Failed to restart safe-socket MockServer [%s] on port %s: %v", name, port, err)
	}

	m := &MockServer{
		server:   server,
		received: make(chan string, 100),
		shutdown: make(chan struct{}),
		port:     port,
		running:  true,
	}

	go m.acceptLoop()
	return m
}

func (m *MockServer) acceptLoop() {
	for {
		conn, err := m.server.Accept()
		if err != nil {
			select {
			case <-m.shutdown:
				return
			default:
				continue
			}
		}
		m.mu.Lock()
		m.conns = append(m.conns, conn)
		m.mu.Unlock()
		go m.handleConnection(conn)
	}
}

func (m *MockServer) handleConnection(conn interfaces.TransportConnection) {
	defer conn.Close()

	for {
		msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Non-blocking send to channel
		select {
		case m.received <- string(msg):
		default:
		}
	}
}

func (m *MockServer) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		close(m.shutdown)
		m.server.Close()
		// Force close all active connections
		for _, c := range m.conns {
			c.Close()
		}
		m.conns = nil
		m.running = false
	}
}

func (m *MockServer) GetLastMessage(timeout time.Duration) (string, bool) {
	select {
	case msg := <-m.received:
		return msg, true
	case <-time.After(timeout):
		return "", false
	}
}

// -----------------------------------------------------------------------------

func TestFullEcosystemResilience(t *testing.T) {
	// 1. Setup Mock Servers using safe-socket
	// Log server uses "tcp" (Raw Framed TCP), Notif server uses "tcp-hello"
	logServer, logHost, logPort := NewMockServer(t, "log-server-mock", "tcp-hello")
	notifServer, notifHost, notifPort := NewMockServer(t, "notif-server-mock", "tcp-hello")

	// 2. Setup Config
	baseConfig, _ := Init("resilience-test", "standalone", "devel", "INFO", false, nil)
	baseConfig.Config.Capabilities["log_server"] = map[string]interface{}{
		"ip":   logHost,
		"port": logPort,
	}
	baseConfig.Config.Capabilities["notif_server"] = map[string]interface{}{
		"ip":   notifHost,
		"port": notifPort,
	}

	// 3. Init Logger with "standard" profile
	opts := BootstrapOptions{
		Name:            "resilience-test",
		LoggerProfile:   "standard",
		InitialLogLevel: utils.LevelInfo,
		ExistingConfig:  baseConfig,
	}

	_, uniLog := InitWithOptions(opts)

	// Since Log Server is "tcp", the first message might be the hello message.
	// We need to consume it if it arrives.
	// Actually, the standard profile DOES send a hello message.

	// 4. Verify baseline connection
	uniLog.Info("Baseline Log")

	// Consume first message for log server (might be hello)
	msg, ok := logServer.GetLastMessage(2 * time.Second)
	if !ok {
		t.Error("Log server did not receive anything")
	}
	// If it was the hello message, the next one should be the log
	if len(msg) < 5 { // simple check to see if it's a capnp hello (usually short) vs log message
		_, ok = logServer.GetLastMessage(2 * time.Second)
		if !ok {
			t.Error("Log server did not receive log message after potential hello")
		}
	}

	uniLog.Warning("Baseline Notif")
	if _, ok := notifServer.GetLastMessage(2 * time.Second); !ok {
		t.Error("Notif server did not receive baseline notification")
	}

	// 5. Simulate CRASH
	logServer.Stop()
	notifServer.Stop()
	t.Log("Servers stopped, logging into the void...")

	for i := 0; i < 3; i++ {
		uniLog.Warning("Message during outage")
	}

	// 6. RESTART Servers
	newLogServer := NewMockServerOnPort(t, "log-server-mock", "tcp", logPort)
	newNotifServer := NewMockServerOnPort(t, "notif-server-mock", "tcp-hello", notifPort)

	t.Log("Servers restarted, waiting for reconnection...")

	// 7. Verify recovery
	recoveredLog := false
	recoveredNotif := false

	for i := 0; i < 20; i++ {
		uniLog.Warning("Recovery Check %d", i)

		if !recoveredLog {
			if _, ok := newLogServer.GetLastMessage(500 * time.Millisecond); ok {
				recoveredLog = true
				t.Log("Log server RECOVERED")
			}
		}

		if !recoveredNotif {
			if _, ok := newNotifServer.GetLastMessage(500 * time.Millisecond); ok {
				recoveredNotif = true
				t.Log("Notif server RECOVERED")
			}
		}

		if recoveredLog && recoveredNotif {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !recoveredLog {
		t.Error("Log server did not recover connection")
	}
	if !recoveredNotif {
		t.Error("Notif server did not recover connection")
	}

	// 8. Cleanup
	uniLog.Close()
	newLogServer.Stop()
	newNotifServer.Stop()
}
