package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Bastien-Antigravity/universal-logger/src/bootstrap"
	"github.com/Bastien-Antigravity/universal-logger/src/utils"
)

// -------------------------------------------------------------------------

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/universal-logger/main.go [scenario]")
		fmt.Println("Scenarios:")
		fmt.Println("  1: Local Development (standalone + devel)")
		fmt.Println("  2: Production Standard (production + standard)")
		fmt.Println("  3: High Throughput (production + high_perf)")
		fmt.Println("  4: Testing (test + minimal)")
		fmt.Println("  5: Monitoring (preprod + notif_logger)")
		fmt.Println("  6: Lock-Free Performance (production + no_lock)")
		return
	}

	scenario := os.Args[1]
	var (
		name          string
		configProfile string
		loggerProfile string
		logLevel      = utils.LevelInfo
		useLocalNotifier = false
	)

	switch scenario {
	case "1":
		// SCENARIO 1: LOCAL DEVELOPMENT
		fmt.Println(">>> Starting Scenario 1: Local Development")
		name = "dev-service"
		configProfile = "standalone"
		loggerProfile = "devel"
		logLevel = utils.LevelDebug

	case "2":
		// SCENARIO 2: PRODUCTION STANDARD
		fmt.Println(">>> Starting Scenario 2: Production Standard")
		name = "prod-api"
		configProfile = "production"
		loggerProfile = "standard"
		logLevel = utils.LevelInfo

	case "3":
		// SCENARIO 3: HIGH THROUGHPUT
		fmt.Println(">>> Starting Scenario 3: High Performance")
		name = "high-load-worker"
		configProfile = "production"
		loggerProfile = "high_perf"
		logLevel = utils.LevelWarning

	case "4":
		// SCENARIO 4: TESTING
		fmt.Println(">>> Starting Scenario 4: Automated Testing")
		name = "test-suite"
		configProfile = "test"
		loggerProfile = "minimal"
		logLevel = utils.LevelError

	case "5":
		// SCENARIO 5: MONITORING & ALERTING
		fmt.Println(">>> Starting Scenario 5: Monitoring focused")
		name = "monitor-svc"
		configProfile = "preprod"
		loggerProfile = "notif_logger"
		logLevel = utils.LevelInfo

	case "6":
		// SCENARIO 6: LOCK-FREE PERFORMANCE
		fmt.Println(">>> Starting Scenario 6: Lock-Free Low Latency")
		name = "latency-critical-app"
		configProfile = "production"
		loggerProfile = "no_lock"
		logLevel = utils.LevelInfo

	default:
		fmt.Printf("Unknown scenario: %s\n", scenario)
		return
	}

	// -------------------------------------------------------------------------
	// EXECUTION

	distConfig, uniLog := bootstrap.Init(name, configProfile, loggerProfile, logLevel, useLocalNotifier)

	uniLog.Info("Facade initialized for scenario %s", scenario)
	uniLog.Warning("This is a warning log from %s", name)

	// -------------------------------------------------------------------------
	// DEMONSTRATION: Dynamic Configuration and Callbacks

	// 1. Register a callback for configuration updates
	distConfig.OnConfigUpdate(func(update map[string]map[string]string) {
		fmt.Println(">>> [Event] Configuration update received via callback!")
	})

	// 2. Manually change the log level
	fmt.Println(">>> [Demo] Switching log level to DEBUG manually")
	uniLog.SetLevel(utils.LevelDebug)
	uniLog.Debug("This debug message is now visible after SetLevel()")

	// 3. Update a configuration value dynamically (e.g., simulated remote update)
	fmt.Println(">>> [Demo] Updating memory configuration dynamically")
	distConfig.SetConfig("system", "status", "maintenance")

	// 4. Verify value retrieval
	status := distConfig.GetConfig("system", "status")
	fmt.Printf(">>> [Verify] Current system status: %s\n", status)

	// 5. Test Local Notifications (if enabled)
	if useLocalNotifier {
		fmt.Println(">>> [Demo] Waiting for a local notification...")
		notifQueue := uniLog.GetNotifQueue()
		
		// In a real app, this would be a background goroutine
		go func() {
			for msg := range notifQueue {
				fmt.Printf(">>> [NOTIF] Received: %s (Tags: %v)\n", msg.Message, msg.Tags)
			}
		}()
		
		// Trigger a notification (assuming the mock server or some logic triggers it)
		// For the demo, we just log something that might trigger a notification
		uniLog.Critical("ALERT: High resource usage detected!") 
	}

	// -------------------------------------------------------------------------

	fmt.Printf("Config Object initialized: %v\n", distConfig.Common.Name)

	uniLog.Info("Walkthrough complete.")
	
	// Wait a moment for background events (callbacks, notifications) to finish
	time.Sleep(200 * time.Millisecond)
}
