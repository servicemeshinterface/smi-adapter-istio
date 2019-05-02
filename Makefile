OPERATOR_IMAGE ?= docker.io/albanc/fili

.PHONY: build
build:
	operator-sdk build $(OPERATOR_IMAGE)

push:
	docker push $(OPERATOR_IMAGE)
