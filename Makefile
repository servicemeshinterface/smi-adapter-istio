OPERATOR_IMAGE ?= docker.io/deislabs/smi-adapter-istio
OPERATOR_SDK_RELEASE_VERSION ?= v0.9.0
OPERATOR_IMAGE_TAG ?= latest

#go options
GO        ?= go
GOFLAGS   :=
TESTS     := .
TESTFLAGS :=

# Required for globs to work correctly
SHELL=/usr/bin/env bash

GIT_COMMIT  ?= $(shell git rev-parse --short HEAD)

HAS_DEP := $(shell command -v dep;)
HAS_OPERATOR_SDK := $(shell command -v operator-sdk)

.PHONY: ci-build
ci-build:
ifndef HAS_OPERATOR_SDK
	# install linux release binary
	curl -JL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_RELEASE_VERSION}/operator-sdk-${OPERATOR_SDK_RELEASE_VERSION}-x86_64-linux-gnu > operator-sdk
	sudo chmod a+x operator-sdk
	sudo mv operator-sdk /usr/local/bin
endif

	operator-sdk build $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG)
	docker tag $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG) $(OPERATOR_IMAGE):$(GIT_COMMIT)

ci-push:
	docker push $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG)
	docker push $(OPERATOR_IMAGE):$(GIT_COMMIT)

.PHONY: build check-env
build: check-env
	operator-sdk build $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG)
	docker tag $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG) $(OPERATOR_IMAGE):$(GIT_COMMIT)

push:
	docker push $(OPERATOR_IMAGE):$(OPERATOR_IMAGE_TAG)
	docker push $(OPERATOR_IMAGE):$(GIT_COMMIT)

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
	$(GO) test $(GOFLAGS) -run $(TESTS) ./cmd/... ./pkg/... $(TESTFLAGS) -v

test-e2e:
	operator-sdk test local ./test/e2e --namespaced-manifest deploy/kubernetes-manifests.yaml --namespace istio-system

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure
