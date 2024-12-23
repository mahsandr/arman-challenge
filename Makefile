# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Build variables
BUILD_DIR ?= build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE_FMT = +%FT%T%z
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
export CGO_ENABLED ?= 1
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

# Project variables

# Dependency versions
GOTESTSUM_VERSION = 1.11.0
GOLANGCI_VERSION = latest
BUF_VERSION = 1.0.0

GOLANG_VERSION = 1.23

# Add the ability to override some variables
# Use with care
-include override.mk

.PHONY: check
check: test-all lint ## Run tests and linters

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

TEST_PKGS ?= ./...
TEST_REPORT_NAME ?= results.xml
.PHONY: test
test: TEST_REPORT ?= main
test: TEST_FORMAT ?= standard-quiet
test: SHELL = /bin/bash
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}/test_results/${TEST_REPORT}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/test_results/${TEST_REPORT}/${TEST_REPORT_NAME} --format ${TEST_FORMAT} -- $(filter-out -v,${GOARGS}) -coverprofile=coverage.out -race -parallel 1 $(if ${TEST_PKGS},${TEST_PKGS},./...)
	@go tool cover -func=coverage.out
	@rm coverage.out



.PHONY: test-all
test-all: ## Run all tests
	@${MAKE} GOARGS="${GOARGS} -run .\* " TEST_REPORT=all test

.PHONY: test-integration
test-integration: ## Run integration tests
	@${MAKE} GOARGS="${GOARGS} -run ^TestIntegration\$$\$$" TEST_REPORT=integration test

.PHONY: test-functional
test-functional: ## Run functional tests
	@${MAKE} GOARGS="${GOARGS} -run ^TestFunctional\$$\$$" TEST_REPORT=functional test

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run --concurrency 1
lint-fix: bin/golangci-lint ## Run linter
	bin/golangci-lint run --concurrency 1 --fix


TAG_PREFIX = v
IMAGE_NAME = localhost:5000/mahsandr/arman-challenge

release-%:
	@sed -e "s/^## \[Unreleased\]$$/## [Unreleased]\\"$$'\n'"\\"$$'\n'"\\"$$'\n'"## [$*] - $$(date +%Y-%m-%d)/g; s|^\[Unreleased\]: \(.*\/compare\/\)\(.*\)...HEAD$$|[Unreleased]: \1${TAG_PREFIX}$*...HEAD\\"$$'\n'"[$*]: \1\2...${TAG_PREFIX}$*|g" CHANGELOG.md > CHANGELOG.md.new
	@mv CHANGELOG.md.new CHANGELOG.md

ifeq (${TAG}, 1)
	git add CHANGELOG.md
	git commit -m 'Prepare release $*'
	git tag -m 'Release $*' ${TAG_PREFIX}$*
ifeq (${PUSH}, 1)
	git push; git push origin ${TAG_PREFIX}$*
endif
endif

	@echo "Version updated to $*!"
ifneq (${PUSH}, 1)
	@echo
	@echo "Review the changes made by this script then execute the following:"
ifneq (${TAG}, 1)
	@echo
	@echo "git add CHANGELOG.md && git commit -m 'Prepare release $*' && git tag -m 'Release $*' ${TAG_PREFIX}$*"
	@echo
	@echo "Finally, push the changes:"
endif
	@echo
	@echo "git push; git push origin ${TAG_PREFIX}$*"
endif

.PHONY: major-beta
major-beta: ## Release a new major beta version
	@$(MAKE) release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .-]' '{print ""$$1+1"."0"."0"-beta-1"}')

.PHONY: minor-beta
minor-beta: ## Release a new minor beta version
	@$(MAKE) release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .-]' '{print ""$$1"."$$2+1"."0"-beta-1"}')

.PHONY: next-beta
nextbeta: ## Release a new beta version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[.-]' '{print $$1"."$$2"."0"-beta-"$$5+1}')

.PHONY: minor
minor: ## Release a new minor version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2+1".0"}')

.PHONY: major
major: ## Release a new major version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1+1".0.0"}')

.PHONY: patch
patch: ## Release a new patch version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2"."$$3+1}')


.PHONY: help
help:
	@echo "Usage:"
	@echo "  make major-beta     - Release a new beta version with major changes"
	@echo "  make minor-beta     - Release a new beta version with minor changes"
	@echo "  make next-beta      - Increase the beta version of the previous release"
	@echo "  make minor          - Release a new minor version"
	@echo "  make major          - Release a new major version"
	@echo "  make patch          - Release a new patch version"
	@echo "  make docker         - Build and push Docker image"
	@echo "  ... (other targets)"

docker: build-release
	docker buildx build --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --build-arg GOPROXYURL="https://goproxy.cn" -t localhost:5000/mahsandr/arman-challenge:${VERSION}  .
	docker push localhost:5000/mahsandr/arman-challenge:${VERSION}

# Variables for database migration and build
GOOSE = $(shell which goose)
GOOSE_CMD = github.com/pressly/goose/v3/cmd/goose
CLICKHOUSE_DB_URL ?= tcp://127.0.0.1:9000
MIGRATIONS_DIR = ./internal/infrastructure/clickhouse/migrations
MAIN_BINARY = $(or ${MAIN_BINARY}, main)

# Targets for database migration and build
.PHONY: all install-goose create-migration migrate build run clean

# Default target
all: install-goose build

# Install goose if not installed
install-goose:
ifndef GOOSE
	@echo "goose not found, installing..."
	go install $(GOOSE_CMD)@latest
else
	@echo "goose is already installed"
endif

# Create a new migration
create-migration:
ifndef NAME
	@echo "Error: Migration name is required. Use 'make create-migration NAME=your_migration_name'"
	exit 1
endif
	@$(GOOSE) -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Run migrations
migrate:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) clickhouse "$(CLICKHOUSE_DB_URL)" up
# Rollback migrations
migrate-down:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) clickhouse "$(CLICKHOUSE_DB_URL)" down
# Build the Go project
build:
	@echo "Building project..."
	go build -o $(MAIN_BINARY) .

# Run the built binary
run: build
	@echo "Running application..."
	./$(MAIN_BINARY)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(MAIN_BINARY)
