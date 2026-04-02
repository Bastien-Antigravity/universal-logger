package main

import (
	"fmt"
	"os"

	"github.com/Bastien-Antigravity/distconf-flexlog/src/factory"
	"github.com/Bastien-Antigravity/distconf-flexlog/src/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/examples/main.go [scenario]")
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
	var params models.MFacadeParams

	switch scenario {
	case "1":
		// SCENARIO 1: LOCAL DEVELOPMENT
		// - Local YAML configuration (standalone)
		// - Console + Text File logging (devel)
		// - Level: DEBUG
		fmt.Println(">>> Starting Scenario 1: Local Development")
		params = models.MFacadeParams{
			ConfigProfile: "standalone",
			AppName:       "dev-service",
			LoggerProfile: "devel",
			LogLevel:      "debug",
		}

	case "2":
		// SCENARIO 2: PRODUCTION STANDARD
		// - Remote Config Server + Authoritative local YAML (production)
		// - Console + Binary File + Network + Notifications (standard)
		// - Level: INFO
		fmt.Println(">>> Starting Scenario 2: Production Standard")
		params = models.MFacadeParams{
			ConfigProfile: "production",
			AppName:       "prod-api",
			LoggerProfile: "standard",
			LogLevel:      "info",
		}

	case "3":
		// SCENARIO 3: HIGH THROUGHPUT
		// - Production config
		// - Network-only Async logging with large buffers (high_perf)
		// - Level: WARNING (only critical stuff to minimize I/O)
		fmt.Println(">>> Starting Scenario 3: High Performance")
		params = models.MFacadeParams{
			ConfigProfile: "production",
			AppName:       "high-load-worker",
			LoggerProfile: "high_perf",
			LogLevel:      "warning",
		}

	case "4":
		// SCENARIO 4: TESTING
		// - Local test defaults (test)
		// - Console-only Async logging (minimal)
		// - Level: ERROR
		fmt.Println(">>> Starting Scenario 4: Automated Testing")
		params = models.MFacadeParams{
			ConfigProfile: "test",
			AppName:       "test-suite",
			LoggerProfile: "minimal",
			LogLevel:      "error",
		}

	case "5":
		// SCENARIO 5: MONITORING & ALERTING
		// - GET-only remote config (preprod)
		// - Focus on notification queuing (notif_logger)
		// - Level: INFO
		fmt.Println(">>> Starting Scenario 5: Monitoring focused")
		params = models.MFacadeParams{
			ConfigProfile: "preprod",
			AppName:       "monitor-svc",
			LoggerProfile: "notif_logger",
			LogLevel:      "info",
		}

	case "6":
		// SCENARIO 6: LOCK-FREE PERFORMANCE
		// - Production config
		// - Fully bitwise/atomic logging without mutexes (no_lock)
		// - Level: INFO
		fmt.Println(">>> Starting Scenario 6: Lock-Free Low Latency")
		params = models.MFacadeParams{
			ConfigProfile: "production",
			AppName:       "latency-critical-app",
			LoggerProfile: "no_lock",
			LogLevel:      "info",
		}

	default:
		fmt.Printf("Unknown scenario: %s\n", scenario)
		return
	}

	// EXECUTION
	app := factory.NewFacade(params)
	defer app.Close()

	app.Info("Facade initialized for scenario %s", scenario)
	app.Warning("This is a warning log from %s", params.AppName)

	cfg := app.GetConfig()
	fmt.Printf("Config Object initialized: %v\n", cfg.Common.Name)

	app.Info("Walkthrough complete.")
}
