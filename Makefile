# =============================================================================
# ESSENTIAL PROCESS: Unified build, test, audit, and clean lifecycle orchestration.
#
# DATA FLOW:
#   1. Builds native Go packages and compiled CGO shared library (libunilog).
#   2. Executes unit and race tests across Go, Rust, and Python bindings.
#   3. Runs parity audit against C header exports across all polyglot targets.
#
# KEY PARAMETERS:
#   - build-lib / shared-lib: Compiles CGO shared library for polyglot runtimes.
#   - audit: Audits FFI symbol parity across language implementations.
# =============================================================================

VERSION := $(strip $(shell cat VERSION.txt))
$(if $(VERSION),,$(error VERSION.txt is missing or empty))

.PHONY: all build build-lib shared-lib core test verify-standalone audit version clean

all: build

version:
	@echo $(VERSION)

verify-standalone:
	@echo "Verifying standalone Go build & tests (GOWORK=off)..."
	GOWORK=off go build ./...
	GOWORK=off go test ./...

build-lib shared-lib core:
	@echo "Building CGO shared library libunilog (version $(VERSION))..."
	@mkdir -p libunilog
	@OS=$$(uname -s); \
	case "$$OS" in \
		Darwin) \
			go build -buildmode=c-shared -o libunilog/libunilog.dylib ./src/cgo_bridge; \
			install_name_tool -id @rpath/libunilog.dylib libunilog/libunilog.dylib 2>/dev/null || true ;; \
		MINGW*|MSYS*|CYGWIN*|Windows_NT) \
			go build -buildmode=c-shared -o libunilog/libunilog.dll ./src/cgo_bridge ;; \
		*) \
			go build -buildmode=c-shared -o libunilog/libunilog.so ./src/cgo_bridge ;; \
	esac

build: build-lib
	@echo "Building repository (version $(VERSION))..."
	go build ./...

test:
	@echo "Running tests (version $(VERSION))..."
	go test -v -race ./...
	pytest unilog/python
	make test -C unilog/cpp
	cd unilog/rust && cargo test

audit:
	@echo "Running polyglot parity audit..."
	python3 scripts/parity-audit.py

clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist build bin *.egg-info target/ logs/ *.log src/bootstrap/logs/ unilog/python/logs/ unilog/cpp/logs/
	rm -f libunilog/*.dylib libunilog/*.so libunilog/*.dll
	rm -f unilog/libunilog/*.dylib unilog/libunilog/*.so unilog/libunilog/*.dll
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
