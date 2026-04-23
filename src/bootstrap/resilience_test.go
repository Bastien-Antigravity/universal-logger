package bootstrap

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// MockServer simulates a remote service (Log or Notif server)
type MockServer struct {
	addr     string
	listener net.Listener
	mu       sync.Mutex
	received chan string
	running  bool
}

func NewMockServer() *MockServer {
	return &MockServer{
		received: make(chan string, 100),
	}
}

func (m *MockServer) Start() (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", "", err
	}
	m.listener = ln
	m.addr = ln.Addr().String()
	m.running = true

	host, port, _ := net.SplitHostPort(m.addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go m.handleConnection(conn)
		}
	}()

	return host, port, nil
}

func (m *MockServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	m.mu.Lock()
	m.mu.Unlock()
	
	b := make([]byte, 4096)
	n, err := conn.Read(b)
	if err == nil && n > 0 {
		msg := string(b[:n])
		// fmt.Printf("!!! MockServer [%s] received: %s\n", m.addr, msg)
		m.received <- msg
	}
}

func (m *MockServer) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listener != nil {
		m.listener.Close()
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
	// 1. Setup Mock Servers
	logServer := NewMockServer()
	logHost, logPort, _ := logServer.Start()
	defer logServer.Stop()

	notifServer := NewMockServer()
	notifHost, notifPort, _ := notifServer.Start()
	defer notifServer.Stop()

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
	defer uniLog.Close()

	// 4. Verify baseline connection
	uniLog.Info("Baseline Log")
	if _, ok := logServer.GetLastMessage(2 * time.Second); !ok {
		t.Error("Log server did not receive baseline log")
	}

	uniLog.Warning("Baseline Notif")
	if _, ok := notifServer.GetLastMessage(2 * time.Second); !ok {
		t.Error("Notif server did not receive baseline notification")
	}

	// 5. Simulate CRASH (Stop servers)
	logServer.Stop()
	notifServer.Stop()
	t.Log("Servers stopped, logging into the void...")

	for i := 0; i < 3; i++ {
		uniLog.Warning("Message during outage")
	}

	// 6. RESTART Servers (on same ports)
	// We need to re-bind manually to the same ports if possible, 
	// but for the test we'll just verify the logger can RE-resolve or re-connect 
	// if we update the config or if it keeps trying the same address.
	// Note: Flexible-logger keeps the IP/Port pointers from the config.
	
	newLogServer := NewMockServer()
	lnLog, err := net.Listen("tcp", "127.0.0.1:"+logPort)
	if err != nil {
		t.Fatalf("Failed to restart LogServer: %v", err)
	}
	newLogServer.listener = lnLog
	go func() {
		for {
			conn, err := lnLog.Accept()
			if err != nil {
				return
			}
			newLogServer.handleConnection(conn)
		}
	}()
	defer lnLog.Close()

	newNotifServer := NewMockServer()
	lnNotif, err := net.Listen("tcp", "127.0.0.1:"+notifPort)
	if err != nil {
		t.Fatalf("Failed to restart NotifServer: %v", err)
	}
	newNotifServer.listener = lnNotif
	go func() {
		for {
			conn, err := lnNotif.Accept()
			if err != nil {
				return
			}
			newNotifServer.handleConnection(conn)
		}
	}()
	defer lnNotif.Close()

	t.Log("Servers restarted, waiting for reconnection...")

	// 7. Verify recovery
	// Standard strategy retries every few seconds (max 2s in some configs)
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
}
