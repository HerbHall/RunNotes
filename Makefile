IMAGE ?= herbhall/runnotes
TAG ?= latest

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
RESET = \033[0m

.PHONY: build-extension install-extension update-extension clean

build-extension:
	docker build --tag=$(IMAGE):$(TAG) .

install-extension:
	docker extension install $(IMAGE):$(TAG)

update-extension:
	docker extension update $(IMAGE):$(TAG)

debug-extension:
	docker extension dev debug $(IMAGE)

reset-extension:
	docker extension dev reset $(IMAGE)

push-extension:
	docker buildx create --name=$(BUILDER) || true
	docker buildx use $(BUILDER)
	docker buildx build --push --platform=linux/amd64,linux/arm64 --tag=$(IMAGE):$(TAG) .

clean:
	docker extension rm $(IMAGE) || true
