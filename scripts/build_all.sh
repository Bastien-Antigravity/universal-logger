#!/bin/bash
set -e

# Configuration
PROJECT_ROOT=$(pwd)
LIB_DIR="${PROJECT_ROOT}/libunilog"
BRIDGE_DIR="./src/cgo_bridge"
IMAGE_NAME="unilog-cross-build"

echo ">>> Cleaning old binaries in ${LIB_DIR}..."
mkdir -p "${LIB_DIR}"
rm -f "${LIB_DIR}"/libunilog.{so,dylib,dll,h}

# -------------------------------------------------------------------------
# 1. BUILD FOR macOS (Local)
# -------------------------------------------------------------------------
echo ">>> Building for macOS (local)..."
go build -buildmode=c-shared -o "${LIB_DIR}/libunilog.dylib" \
    "${BRIDGE_DIR}/initialize.go" \
    "${BRIDGE_DIR}/config.go" \
    "${BRIDGE_DIR}/logger.go"

# -------------------------------------------------------------------------
# 2. CHECK DOCKER DAEMON
# -------------------------------------------------------------------------
if ! docker info > /dev/null 2>&1; then
    echo "ERROR: Docker daemon is not running. Cross-compilation for Linux and Windows requires Docker."
    echo "Please start Docker Desktop on your Mac and try again."
    # We continue so that macOS build is still usable
    exit 0
fi

if [[ "$(docker images -q ${IMAGE_NAME} 2> /dev/null)" == "" ]]; then
    echo ">>> Building Docker cross-compilation image..."
    docker build -t ${IMAGE_NAME} -f scripts/Dockerfile.cross scripts/
else
    echo ">>> Docker image ${IMAGE_NAME} already exists. Skipping build."
fi


# -------------------------------------------------------------------------
# 3. BUILD FOR LINUX (x86_64) via Docker
# -------------------------------------------------------------------------
echo ">>> Building for Linux x86_64..."
docker run --rm -v "${PROJECT_ROOT}:/app" ${IMAGE_NAME} \
    "CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-gnu-gcc \
    go build -buildmode=c-shared -o libunilog/libunilog.so \
    ${BRIDGE_DIR}/initialize.go ${BRIDGE_DIR}/config.go ${BRIDGE_DIR}/logger.go"

# -------------------------------------------------------------------------
# 4. BUILD FOR WINDOWS (x86_64) via Docker
# -------------------------------------------------------------------------
echo ">>> Building for Windows x86_64..."
docker run --rm -v "${PROJECT_ROOT}:/app" ${IMAGE_NAME} \
    "CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
    go build -buildmode=c-shared -o libunilog/libunilog.dll \
    ${BRIDGE_DIR}/initialize.go ${BRIDGE_DIR}/config.go ${BRIDGE_DIR}/logger.go"

echo ">>> All builds complete! Artifacts are in: ${LIB_DIR}"
ls -lh "${LIB_DIR}"/libunilog.*


echo ">>> All builds complete!"
ls -lh "${PYTHON_DIR}"/libunilog.*
