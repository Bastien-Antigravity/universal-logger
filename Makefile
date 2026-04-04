# Root Makefile for Universal Logger
# Orchestrates builds for Go (Core), C++, Python, and Rust.

.PHONY: all core cpp python rust clean help

# Platform Detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    LIB_EXT := dylib
else
    LIB_EXT := so
endif

LIB_DIR := $(PWD)/libunilog
CORE_SRC := $(wildcard src/cgo_bridge/*.go)

all: core cpp rust python

help:
	@echo "Universal Logger Build System"
	@echo "Usage:"
	@echo "  make all      - Build everything (Core + all clients)"
	@echo "  make core     - Build Go shared library (CGO)"
	@echo "  make cpp      - Build C++ example"
	@echo "  make rust     - Build Rust example"
	@echo "  make python   - Run Python tests (requires Core)"
	@echo "  make clean    - Remove all build artifacts"

core: $(LIB_DIR)/libunilog.$(LIB_EXT)

$(LIB_DIR)/libunilog.$(LIB_EXT): $(CORE_SRC)
	@echo ">>> Building Go Shared Library (CGO)..."
	@mkdir -p $(LIB_DIR)
	go build -buildmode=c-shared -o $(LIB_DIR)/libunilog.$(LIB_EXT) $(CORE_SRC)
	@if [ "$(UNAME_S)" = "Darwin" ]; then \
		install_name_tool -id "@rpath/libunilog.dylib" $(LIB_DIR)/libunilog.dylib; \
	fi

cpp: core
	@echo ">>> Building C++ Client..."
	$(MAKE) -C cpp

rust: core
	@echo ">>> Building Rust Client (Library + Demo)..."
	@# Ensure Rust knows where to find the library at link time
	cd rust && RUSTFLAGS="-L $(LIB_DIR)" cargo build --example demo

python: core
	@echo ">>> Running Python Tests..."
	@# Ensure Python can find the library
	export DYLD_LIBRARY_PATH=$(LIB_DIR):$$DYLD_LIBRARY_PATH && \
	export LD_LIBRARY_PATH=$(LIB_DIR):$$LD_LIBRARY_PATH && \
	export PYTHONPATH=$(PWD)/python:$$PYTHONPATH && \
	python3 python/test_unilog.py && \
	python3 python/test_unified_callback.py && \
	python3 python/test_async_logging.py

clean:
	@echo ">>> Cleaning all build artifacts..."
	rm -rf $(LIB_DIR)
	$(MAKE) -C cpp clean
	cd rust && cargo clean
