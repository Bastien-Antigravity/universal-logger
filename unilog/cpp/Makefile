# Makefile for Universal Logger C++ Example
# Builds the C++ client linking against the Go shared library.

# Platform detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    LIB_EXT := dylib
else
    LIB_EXT := so
endif

# Paths
LIB_DIR := ../libunilog
LIB_NAME := unilog
TARGET := unilog_cpp
TEST_TARGET := test_unilog
CORE_SRC := $(wildcard ../../src/cgo_bridge/*.go)

# Compilation
CXX := g++
CXXFLAGS := -std=c++11 -Wall -Wno-unused-function -I$(LIB_DIR)
LDFLAGS := -L$(LIB_DIR) -l$(LIB_NAME)

# Ensure the library can be found at runtime on macOS
ifeq ($(UNAME_S),Darwin)
    LDFLAGS += -Wl,-rpath,@loader_path/$(LIB_DIR)
endif

all: $(TARGET) $(TEST_TARGET)

# Rule to build the Go shared library and header
$(LIB_DIR)/libunilog.$(LIB_EXT) $(LIB_DIR)/libunilog.h: $(CORE_SRC)
	@echo ">>> Building Go shared library (CGO)..."
	@mkdir -p $(LIB_DIR)
	go build -buildmode=c-shared -o $(LIB_DIR)/libunilog.$(LIB_EXT) $(CORE_SRC)

$(TARGET): main.cpp $(LIB_DIR)/libunilog.h
	$(CXX) $(CXXFLAGS) main.cpp -o $(TARGET) $(LDFLAGS)

$(TEST_TARGET): test_unilog.cpp $(LIB_DIR)/libunilog.h
	$(CXX) $(CXXFLAGS) test_unilog.cpp -o $(TEST_TARGET) $(LDFLAGS)

clean:
	rm -f $(TARGET) $(TEST_TARGET)
	rm -rf *.dSYM

test: $(TEST_TARGET)
	export DYLD_LIBRARY_PATH=$(LIB_DIR):$$DYLD_LIBRARY_PATH && \
	export LD_LIBRARY_PATH=$(LIB_DIR):$$LD_LIBRARY_PATH && \
	./$(TEST_TARGET)

run: all
	./$(TARGET)
