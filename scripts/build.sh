#!/usr/bin/env bash

STATIC=${1:-false}
PATH=${PATH}:/usr/local/go/bin

GOOS=linux
GOARCH=amd64
BINARY_DIR=bin
BINARY=${BINARY_DIR}/harmony-tui

source $(go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh

if [ "$(uname -s)" == "Darwin" ]; then
   GOOS=darwin
fi

mkdir -p ${BINARY_DIR}

version=$(git rev-list --count HEAD)
commit=$(git describe --always --long --dirty)
built_at=$(date +%FT%T%z)
built_by=${USER}

export GO111MODULE=on
export ldflags="-X main.version=v${version} -X main.commit=${commit} -X main.builtAt=${built_at} -X main.builtBy=${built_by}"

if [ "$STATIC" == "true" ]; then
   env GOOS=${GOOS} GOARCH=amd64 go build -v -ldflags="${ldflags} -linkmode external -extldflags -static" -o $BINARY
else
   env GOOS=${GOOS} GOARCH=amd64 go build -v -ldflags="${ldflags}" -o $BINARY
fi
