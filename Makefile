#   Copyright 2020 Docker Compose CLI authors

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

export DOCKER_BUILDKIT=1

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	MOBY_DOCKER=/usr/bin/docker
endif
ifeq ($(UNAME_S),Darwin)
	MOBY_DOCKER=/Applications/Docker.app/Contents/Resources/bin/docker
endif

BINARY_FOLDER=$(shell pwd)/bin
GIT_TAG?=$(shell git describe --tags --match "v[0-9]*")
TEST_FLAGS?=
E2E_TEST?=
ifeq ($(E2E_TEST),)
else
	TEST_FLAGS=-run $(E2E_TEST)
endif

all: cli

cli: ## Compile the cli
	@docker build . --target cli \
	--platform local \
	--build-arg BUILD_TAGS=e2e,kube \
	--build-arg GIT_TAG=$(GIT_TAG) \
	--output ./bin

e2e-win-ci: ## Run end to end local tests on Windows CI, no Docker for Linux containers available ATM. Set E2E_TEST=TestName to run a single test
	go test -count=1 -v $(TEST_FLAGS) ./local/e2e/cli-only

e2e-ecs: ## Run End to end ECS tests. Set E2E_TEST=TestName to run a single test
	go test -timeout 30m -count=1 -v $(TEST_FLAGS) ./ecs/e2e/ecs ./ecs/e2e/ecs-local

cross: ## Compile the CLI for linux, darwin and windows
	@docker build . --target cross \
	--build-arg BUILD_TAGS \
	--build-arg GIT_TAG=$(GIT_TAG) \
	--output ./bin \

test: ## Run unit tests
	@docker build --progress=plain . \
	--build-arg BUILD_TAGS=kube \
	--build-arg GIT_TAG=$(GIT_TAG) \
	--target test

cache-clear: ## Clear the builder cache
	@docker builder prune --force --filter type=exec.cachemount --filter=unused-for=24h

lint: ## run linter(s)
	@docker build . \
	--build-arg BUILD_TAGS=kube,e2e \
	--build-arg GIT_TAG=$(GIT_TAG) \
	--target lint

check-dependencies: ## check dependency updates
	go list -u -m -f '{{if not .Indirect}}{{if .Update}}{{.}}{{end}}{{end}}' all

import-restrictions: ## run import-restrictions script
	@docker build . \
	--target import-restrictions

validate-headers: ## Check license header for all files
	@docker build . --target check-license-headers

go-mod-tidy: ## Run go mod tidy in a container and output resulting go.mod and go.sum
	@docker build . --target go-mod-tidy --output .

validate-go-mod: ## Validate go.mod and go.sum are up-to-date
	@docker build . --target check-go-mod

validate: validate-go-mod validate-headers ## Validate sources

pre-commit: validate import-restrictions check-dependencies lint cli test e2e-local

build-ecs-search-sidecar:  ## build ecs search sidecar image locally and tag it with make build-ecs-search-sidecar tag=0.1
	docker buildx build --platform linux/amd64,linux/arm64 -t docker/ecs-searchdomain-sidecar:$(tag) ecs/resolv

publish-ecs-search-sidecar: build-ecs-search-sidecar ## build & publish ecs search sidecar image with make publish-ecs-search-sidecar tag=0.1
	docker pull docker/ecs-searchdomain-sidecar:$(tag) && echo "Failure: Tag already exists" || docker buildx build --push --platform linux/amd64,linux/arm64 -t docker/ecs-searchdomain-sidecar:$(tag) ecs/resolv

build-ecs-secrets-sidecar:  ## build ecs secrets sidecar image locally and tag it with make build-ecs-secrets-sidecar tag=0.1
	docker buildx build --platform linux/amd64,linux/arm64 -t docker/ecs-secrets-sidecar:$(tag) ecs/secrets

publish-ecs-secrets-sidecar: build-ecs-secrets-sidecar ## build & publish ecs secrets sidecar image with make publish-ecs-secrets-sidecar tag=0.1
	docker pull docker/ecs-secrets-sidecar:$(tag) && echo "Failure: Tag already exists" || docker buildx build --push --platform linux/amd64,linux/arm64 -t docker/ecs-secrets-sidecar:$(tag) ecs/secrets

help: ## Show help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

FORCE:

.PHONY: all validate cli cross test cache-clear lint check-dependencies help go-mod-tidy
