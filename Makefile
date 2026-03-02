VERSION := $(shell cat VERSION)
IMAGE ?= herbhall/runnotes
TAG ?= $(VERSION)

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
RESET = \033[0m

.PHONY: test lint build-backend dev-backend cross-check build-extension install-extension update-extension clean fe-install fe-build fe-lint fe-typecheck fe-test

test: ## Run Go tests
	go test ./...

lint: ## Run Go linter
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1 run ./...

build-backend: ## Build backend binary (local arch)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/backend ./cmd/backend

dev-backend: ## Run backend in dev mode (TCP :3001)
	ENV_MODE=dev go run ./cmd/backend

cross-check: ## Verify cross-compilation for linux/amd64 and linux/arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dev/null ./cmd/backend
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/backend

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

release: ## Build, push multi-arch image with version tag and latest
	docker buildx create --name=$(BUILDER) || true
	docker buildx use $(BUILDER)
	docker buildx build --push --platform=linux/amd64,linux/arm64 --tag=$(IMAGE):$(TAG) --tag=$(IMAGE):latest .

clean:
	docker extension rm $(IMAGE) || true

## Frontend
fe-install: ## Install frontend dependencies
	cd ui && npm ci

fe-build: ## Build frontend
	cd ui && npm run build

fe-lint: ## Lint frontend
	cd ui && npm run lint

fe-typecheck: ## TypeScript type check
	cd ui && npx tsc --noEmit

fe-test: ## Run frontend tests
	cd ui && npm run test
