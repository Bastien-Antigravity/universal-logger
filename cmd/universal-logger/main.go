package main

import (
	"fmt"
	"os"

	"github.com/Bastien-Antigravity/universal-logger/src/bootstrap"
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
	var name, configProfile, loggerProfile, logLevel string

	switch scenario {
	case "1":
		// SCENARIO 1: LOCAL DEVELOPMENT
		fmt.Println(">>> Starting Scenario 1: Local Development")
		name = "dev-service"
		configProfile = "standalone"
		loggerProfile = "devel"
		logLevel = "debug"

	case "2":
		// SCENARIO 2: PRODUCTION STANDARD
		fmt.Println(">>> Starting Scenario 2: Production Standard")
		name = "prod-api"
		configProfile = "production"
		loggerProfile = "standard"
		logLevel = "info"

	case "3":
		// SCENARIO 3: HIGH THROUGHPUT
		fmt.Println(">>> Starting Scenario 3: High Performance")
		name = "high-load-worker"
		configProfile = "production"
		loggerProfile = "high_perf"
		logLevel = "warning"

	case "4":
		// SCENARIO 4: TESTING
		fmt.Println(">>> Starting Scenario 4: Automated Testing")
		name = "test-suite"
		configProfile = "test"
		loggerProfile = "minimal"
		logLevel = "error"

	case "5":
		// SCENARIO 5: MONITORING & ALERTING
		fmt.Println(">>> Starting Scenario 5: Monitoring focused")
		name = "monitor-svc"
		configProfile = "preprod"
		loggerProfile = "notif_logger"
		logLevel = "info"

	case "6":
		// SCENARIO 6: LOCK-FREE PERFORMANCE
		fmt.Println(">>> Starting Scenario 6: Lock-Free Low Latency")
		name = "latency-critical-app"
		configProfile = "production"
		loggerProfile = "no_lock"
		logLevel = "info"

	default:
		fmt.Printf("Unknown scenario: %s\n", scenario)
		return
	}

	// EXECUTION
	distConfig, logSvc := bootstrap.Init_UniLogger(name, configProfile, loggerProfile, logLevel)
	defer logSvc.Close()

	logSvc.Info("Facade initialized for scenario %s", scenario)
	logSvc.Warning("This is a warning log from %s", name)

	fmt.Printf("Config Object initialized: %v\n", distConfig.Common.Name)

	logSvc.Info("Walkthrough complete.")
}
