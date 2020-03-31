OPERATOR_IMAGE ?= docker.io/servicemeshinterface/smi-adapter-istio
OPERATOR_SDK_RELEASE_VERSION ?= v0.9.0
OPERATOR_IMAGE_TAG ?= latest

#go options
GO        ?= go
GOFLAGS   :=
TESTS     := .
TESTFLAGS :=

KIND_VERSION=v0.5.1
HOST_OS := $(shell uname -s)
ifeq ($(HOST_OS), Darwin)
OS_ARCH = darwin-amd64
else
ifeq ($(HOST_OS), Linux)
OS_ARCH = linux-amd64
else
$(error Unsupported Host OS)
endif
endif

# Required for globs to work correctly
SHELL=/usr/bin/env bash

GIT_COMMIT  ?= $(shell git rev-parse --short HEAD)

HAS_DEP := $(shell command -v dep;)
HAS_OPERATOR_SDK := $(shell command -v operator-sdk)
HAS_KIND := $(shell command -v kind)

KIND_KUBECONFIG = $(shell kind get kubeconfig-path --name=local-kind)

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
	KUBECONFIG=$(KIND_KUBECONFIG) \
		operator-sdk test local ./test/e2e --namespaced-manifest deploy/operator-and-rbac.yaml --namespace istio-system

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure

create-kindcluster:
ifndef HAS_KIND
	@echo "installing kind"
	@curl -fsSLo kind "https://github.com/kubernetes-sigs/kind/releases/download/$(KIND_VERSION)/kind-$(OS_ARCH)"
	@chmod +x kind
	@mv kind /usr/local/bin/kind
endif
	kind create cluster --name local-kind  --image kindest/node:v1.15.3

install-istio:
	curl -fsSL https://git.io/getLatestIstio | ISTIO_VERSION=1.1.6 sh -
	ls istio-1.1.6/install/kubernetes/helm/istio-init/files/crd*yaml | \
		xargs -I{} kubectl apply -f {} --kubeconfig=$(KIND_KUBECONFIG)
	kubectl apply -f istio-1.1.6/install/kubernetes/istio-demo-auth.yaml \
		--kubeconfig=$(KIND_KUBECONFIG)

install-smi-crds:
	kubectl apply -f https://raw.githubusercontent.com/servicemeshinterface/smi-adapter-istio/master/deploy/crds/crds.yaml \
		--kubeconfig=$(KIND_KUBECONFIG)

clean-kind:
	kind delete cluster --name local-kind
