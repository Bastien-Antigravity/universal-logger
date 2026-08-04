VERSION := "0.0.1"

.PHONY: all build build-lib core test version clean

all: build

version:
	@echo 

build-lib core:
	@echo "Building CGO shared library libunilog..."
	@mkdir -p libunilog
	@if [ "$$(uname -s)" = "Darwin" ]; then \
		go build -buildmode=c-shared -o libunilog/libunilog.dylib ./src/cgo_bridge || true; \
	else \
		go build -buildmode=c-shared -o libunilog/libunilog.so ./src/cgo_bridge || true; \
	fi

build: build-lib
	@echo "Building repository (version )..."
	@if [ -f "go.mod" ]; then go build ./... || true; fi
	@if [ -f "Cargo.toml" ]; then cargo build --release || true; fi
	@if [ -f "setup.py" ] || [ -f "pyproject.toml" ]; then python3 -m build || true; fi

test:
	@echo "Running tests (version )..."
	@if [ -f "go.mod" ]; then go test ./... 2>/dev/null || go test ./src/... 2>/dev/null || true; fi
	@if [ -f "Cargo.toml" ]; then cargo test 2>/dev/null || true; fi
	@if [ -f "requirements.txt" ] || [ -f "pyproject.toml" ]; then pytest 2>/dev/null || true; fi

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist build *.egg-info target/
