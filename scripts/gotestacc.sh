#!/usr/bin/env bash

set -e

echo "==> Running acceptance testing..."

PACKAGES="$(find . -name '*.go' | xargs -I{} dirname {} |  cut -f2 -d/ | sort -u | grep -Ev '(^\.$|.git|vendor|bin)' | sed -e 's!^!./!' -e 's!$!/...!')"
TF_ACC=1 go test -cover -tags=test ${PACKAGES} -v -timeout 120m
