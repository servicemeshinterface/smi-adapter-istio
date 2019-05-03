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
