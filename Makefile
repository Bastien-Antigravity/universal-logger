# Root Makefile for Universal Logger
# Orchestrates builds for Go (Core), C++, Python, and Rust.

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

.PHONY: all core cpp python rust clean help audit

# Platform Detection (Matches distributed-config logic using go env)
OS=$(shell go env GOOS)
ifeq ($(OS),windows)
	LIB_EXT=dll
else ifeq ($(OS),darwin)
	LIB_EXT=dylib
else
	LIB_EXT=so
endif

LIB_DIR := $(PWD)/unilog/libunilog
PYTHON_PKG_DIR := $(PWD)/unilog/python/unilog
CORE_SRC := $(wildcard src/cgo_bridge/*.go)

all: core cpp rust python audit

build: core

help:
	@echo "Universal Logger Build System"
	@echo "Usage:"
	@echo "  make all      - Build everything (Core + all clients)"
	@echo "  make core     - Build Go shared library (CGO) & bundle in Python"
	@echo "  make cpp      - Build C++ example"
	@echo "  make rust     - Build Rust example"
	@echo "  make python   - Run Python tests (requires Core)"
	@echo "  make audit    - Run polyglot parity audit"
	@echo "  make clean    - Remove all build artifacts"

audit:
	@echo ">>> Running Polyglot Parity Audit..."
	@python3 scripts/parity-audit.py

core: $(LIB_DIR)/libunilog.$(LIB_EXT)

$(LIB_DIR)/libunilog.$(LIB_EXT): $(CORE_SRC)
	@echo ">>> Building Go Shared Library (CGO)..."
	@mkdir -p $(LIB_DIR)
	$(GOBUILD) -buildmode=c-shared -o $(LIB_DIR)/libunilog.$(LIB_EXT) ./src/cgo_bridge
	@echo ">>> Bundling library in Python package..."
	@mkdir -p $(PYTHON_PKG_DIR)
	cp $(LIB_DIR)/libunilog.$(LIB_EXT) $(PYTHON_PKG_DIR)/
	@if [ "$(LIB_EXT)" = "dylib" ]; then \
		install_name_tool -id "@rpath/libunilog.dylib" $(LIB_DIR)/libunilog.dylib; \
		install_name_tool -id "@rpath/libunilog.dylib" $(PYTHON_PKG_DIR)/libunilog.dylib; \
	fi

cpp: core
	@echo ">>> Building C++ Client..."
	$(MAKE) -C unilog/cpp

rust: core
	@echo ">>> Building Rust Client (Library + Demo)..."
	@# Ensure Rust knows where to find the library at link time
	cd unilog/rust && RUSTFLAGS="-L $(LIB_DIR)" cargo build --example demo

python: core
	@echo ">>> Running Python Tests..."
	@# Ensure Python can find the library
	export DYLD_LIBRARY_PATH=$(LIB_DIR):$$DYLD_LIBRARY_PATH && \
	export LD_LIBRARY_PATH=$(LIB_DIR):$$LD_LIBRARY_PATH && \
	export PYTHONPATH=$(PWD)/unilog/python:$$PYTHONPATH && \
	cd unilog/python && \
	python3 test_unilog.py && \
	python3 test_unified_callback.py && \
	python3 test_async_logging.py

clean:
	@echo ">>> Cleaning all build artifacts..."
	$(GOCLEAN)
	rm -rf $(LIB_DIR)
	rm -f $(PYTHON_PKG_DIR)/libunilog.*
	$(MAKE) -C unilog/cpp clean
	cd unilog/rust && cargo clean
