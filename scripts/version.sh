#!/usr/bin/env bash

if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
    DIRTY="-dirty"
fi

GIT_TAG=$(git tag -l --contains HEAD | head -n 1)

if [ -n "$VERSION" ]; then
    VERSION="$VERSION${DIRTY}"
elif [ -n "$GIT_TAG" ]; then
    VERSION="$GIT_TAG${DIRTY}"
else
    COMMIT=$(git rev-parse --short HEAD)
    VERSION="${COMMIT}${DIRTY}"
fi
