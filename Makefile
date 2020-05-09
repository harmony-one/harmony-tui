TOP:=$(realpath ..)
SHELL := /bin/bash
version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@harmony.one

ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}

uname := $(shell uname)
BINARY_DIR=$(PWD)/bin
BINARY_NAME=harmony-tui

all:
	export GO111MODULE=on
	export ldflags="-X main.version=v${version} -X main.commit=${commit} -X main.builtAt=${built_at} -X main.builtBy=${built_by}"
	mkdir -p ${BINARY_DIR}
	source $(TOP)/harmony/scripts/setup_bls_build_flags.sh && go build -ldflags="$(ldflags)" -o ${BINARY_DIR}/$(BINARY_NAME)

clean:
	rm -rf $(BINARY_DIR)

run: build
	$(BINARY_DIR)/$(BINARY_NAME)

static:
	make -C $(TOP)/mcl
	make -C $(TOP)/bls minimised_static BLS_SWAP_G=1
	HMY_PATH=$(realpath ..) source $(TOP)/harmony/scripts/setup_bls_build_flags.sh && go build -ldflags="$(ldflags) -w -extldflags \"-static\"" -o ${BINARY_DIR}/$(BINARY_NAME)
