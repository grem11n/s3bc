.PHONY: help fmt mod lint test build \
	docker-build-dev docker-lint docker-test

COMMIT ?= $(shell git rev-parse --short=6 HEAD)
NAME ?= s3bc
USER ?= $(shell git config user.email)
DATETIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
DOCKERFILE_DIR = build
DOCKERFILE_DEV = Dockerfile.dev
version ?= $(COMMIT)

help: ## Show help
	@IFS=$$'\n' ; \
    help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
    for help_line in $${help_lines[@]}; do \
        IFS=$$'#' ; \
        help_split=($$help_line) ; \
        help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
        help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
        printf "%-30s %s\n" $$help_command $$help_info ; \
    done

mod: ## Run mod tidy and mod vendor
	@go mod tidy
	@go mod vendor

fmt: ## Autoformat
	@gofmt -w $$(find . -name '*.go' | grep -v vendor)

lint: ## Run linter
	@golangci-lint run -v --timeout=5m --modules-download-mode=vendor -E gosec -E revive -E goconst -E misspell -E whitespace ./...

test: ## Run tests
	@go test -v ./...

build: ## Build the binary (intended to use with Docker)
	@go build -mod=vendor -ldflags="-s -w -X github.com/grem11n/s3bc/version.Version=$(version) -X github.com/grem11n/s3bc/version.Commit=$(COMMIT) -X github.com/grem11n/s3bc/version.Date=$(DATETIME) -X github.com/grem11n/s3bc/version.BuiltBy=$(USER)" -o bin/s3bc
	@echo "Built a new version of s3bc: $(version). Commit hash: $(COMMIT)"


docker-lint: ## Run linter in a Docker container. Requires a built image
	docker run --rm --entrypoint="/usr/bin/make" -e GO111MODULE=on -v $(PWD):/app -w /app $(NAME):$(COMMIT) lint

docker-test: ## Run tests in a Docker container. Requires a built image
	docker run --rm --entrypoint="/usr/bin/make" -e GO111MODULE=on -v $(PWD):/app -w /app $(NAME):$(COMMIT) test

docker-build-dev: ## Builds a Dev Docker image for BCM CLI
	@docker build --no-cache -t $(NAME):$(COMMIT) -t $(NAME):dev -f $(DOCKERFILE_DIR)/$(DOCKERFILE_DEV) --progress=plain .
