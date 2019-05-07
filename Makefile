# We don't have a default name yet.
# OPERATOR_IMAGE ?= docker.io/servicemeshinterface/smi-adapter-istio

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
