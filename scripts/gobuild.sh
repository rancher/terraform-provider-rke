#!/usr/bin/env bash

set -e

source $(dirname $0)/version.sh

echo "==> Building code binaries version ${VERSION} ..."

declare -A OS_ARCH_ARG

OS_PLATFORM_ARG=(linux windows darwin freebsd openbsd)
OS_ARCH_ARG[linux]="amd64 arm arm64"
OS_ARCH_ARG[windows]="386 amd64"
OS_ARCH_ARG[darwin]="amd64"
OS_ARCH_ARG[freebsd]="386 amd64 arm"
OS_ARCH_ARG[openbsd]="386 amd64"

BIN_NAME="terraform-provider-rke"
BUILD_DIR=$(dirname $0)"/../build/bin"


CGO_ENABLED=0 go build -ldflags="-w -s -X main.VERSION=$VERSION -extldflags -static" -o bin/${BIN_NAME}

if [ -n "$CROSS" ]; then
    rm -rf ${BUILD_DIR}
    mkdir -p ${BUILD_DIR}
    for OS in ${OS_PLATFORM_ARG[@]}; do
        for ARCH in ${OS_ARCH_ARG[${OS}]}; do
            OUTPUT_BIN="${BUILD_DIR}/${BIN_NAME}_${OS}_${ARCH}"
            if test "$OS" = "windows"; then
                OUTPUT_BIN="${OUTPUT_BIN}.exe"
            fi
            echo "Building ${BIN_NAME}_${OS}_${ARCH} ..."
            GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build \
                  -ldflags="-w -X main.VERSION=$VERSION" \
                  -o ${OUTPUT_BIN} ./
        done
    done
fi