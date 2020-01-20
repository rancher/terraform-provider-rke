#!/usr/bin/env bash

set -e

echo "==> Building code binaries..."

source $(dirname $0)/version.sh

declare -A OS_ARCH_ARG

OS_PLATFORM_ARG=(linux windows darwin)
OS_ARCH_ARG[linux]="amd64 arm"
OS_ARCH_ARG[windows]="386 amd64"
OS_ARCH_ARG[darwin]="amd64"

BIN_NAME="terraform-provider-rke"
BUILD_DIR=$(dirname $0)"/../build/bin"


CGO_ENABLED=0 go build -ldflags="-w -s -X main.VERSION=$VERSION -extldflags -static" -o bin/${BIN_NAME}

if [ -n "$CROSS" ]; then
    rm -rf ${BUILD_DIR}
    mkdir -p ${BUILD_DIR}
    for OS in ${OS_PLATFORM_ARG[@]}; do
        for ARCH in ${OS_ARCH_ARG[${OS}]}; do
            OUTPUT_BIN="${BUILD_DIR}/${BIN_NAME}_$OS-$ARCH"
            if test "$OS" = "windows"; then
                OUTPUT_BIN="${OUTPUT_BIN}.exe"
            fi
            echo "Building binary for $OS/$ARCH..."
            GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build \
                  -ldflags="-w -X main.VERSION=$VERSION" \
                  -o ${OUTPUT_BIN} ./
        done
    done
fi