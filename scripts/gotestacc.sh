#!/usr/bin/env bash

set -e

echo "==> Running acceptance testing..."

cleanup()
{
    $(dirname $0)/cleanup_testacc.sh
}
trap cleanup EXIT TERM

PACKAGES="$(find . -name '*.go' | xargs -I{} dirname {} |  cut -f2 -d/ | sort -u | grep -Ev '(^\.$|.git|vendor|bin)' | sed -e 's!^!./!' -e 's!$!/...!')"
TF_ACC=1 go test -cover -tags=test ${PACKAGES} -v -timeout 120m
