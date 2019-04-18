#!/usr/bin/env bash

set -e

mkdir -p bin/ 2>/dev/null

for GOOS in $OS; do
    for GOARCH in $ARCH; do
        arch="$GOOS-$GOARCH"
        binary="terraform-provider-rke_v${CURRENT_VERSION}"
        if [ "$GOOS" = "windows" ]; then
          binary="${binary}.exe"
        fi
        echo "Building $binary $arch"
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
            go build \
                -ldflags "$BUILD_LDFLAGS" \
                -o bin/$binary \
                main.go
        if [ -n "$ARCHIVE" ]; then
            (cd bin/; zip -r "terraform-provider-rke_${CURRENT_VERSION}_$arch.zip" $binary)
            rm -f bin/$binary
        fi
    done
done
