# We don't have a default name yet.
# OPERATOR_IMAGE ?= docker.io/servicemeshinterface/smi-adapter-istio


#go options
GO        ?= go
GOFLAGS   :=
PKG       := $(shell glide novendor)
TESTS     := .
TESTFLAGS :=

# Required for globs to work correctly
SHELL=/usr/bin/env bash

.PHONY: build check-env
build: check-env
	operator-sdk build $(OPERATOR_IMAGE)

push:
	docker push $(OPERATOR_IMAGE)

check-env:
ifndef OPERATOR_IMAGE
	$(error Environment variable OPERATOR_IMAGE is undefined)
endif

GOFORMAT_FILES := $(shell find . -name '*.go' | grep -v '\./vendor/')

.PHONY: format-go-code
## Formats any go file that differs from gofmt's style
format-go-code:
	@gofmt -s -l -w ${GOFORMAT_FILES}

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	$(GO) test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS) -v

HAS_GLIDE := $(shell command -v glide;)
HAS_DEP := $(shell command -v dep;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure
