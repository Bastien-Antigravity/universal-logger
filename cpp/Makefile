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
CORE_SRC := $(wildcard ../src/cgo_bridge/*.go)

# Compilation
CXX := g++
CXXFLAGS := -std=c++11 -Wall -I$(LIB_DIR)
LDFLAGS := -L$(LIB_DIR) -l$(LIB_NAME)

# Ensure the library can be found at runtime on macOS
ifeq ($(UNAME_S),Darwin)
    LDFLAGS += -Wl,-rpath,@loader_path/$(LIB_DIR)
endif

all: $(TARGET)

# Rule to build the Go shared library and header
$(LIB_DIR)/libunilog.$(LIB_EXT) $(LIB_DIR)/libunilog.h: $(CORE_SRC)
	@echo ">>> Building Go shared library (CGO)..."
	@mkdir -p $(LIB_DIR)
	go build -buildmode=c-shared -o $(LIB_DIR)/libunilog.$(LIB_EXT) $(CORE_SRC)

$(TARGET): main.cpp $(LIB_DIR)/libunilog.h
	$(CXX) $(CXXFLAGS) main.cpp -o $(TARGET) $(LDFLAGS)

clean:
	rm -f $(TARGET)

run: all
	./$(TARGET)
