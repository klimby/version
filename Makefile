#!/usr/bin/make
# Makefile readme (ru): <http://linux.yaroslavl.ru/docs/prog/gnu_make_3-79_russian_manual.html>
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>
SHELL = /bin/sh

# Get current git tag
GIT_VERSION_TAG := $(shell git describe --tags $$(git rev-list --tags --max-count=1))

# Package version - git tag without 'v' prefix
PACKAGE_VERSION := $(patsubst v%,%,$(GIT_VERSION_TAG))

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: build
build: ## Build to bin folder
	$(eval V := $(or $(VERSION),$(PACKAGE_VERSION)))
	@go build -ldflags "-s -w -X main.version=$(V)" -o ./bin/version github.com/klimby/version
	@sudo chmod +x ./bin/version
	@echo "Build created v$(V)"
	@./bin/version --version

.PHONY: build-self
build-self: ## Build to root folder for use in this project
	$(eval V := $(or $(VERSION),$(PACKAGE_VERSION)))
	go build  -ldflags "-s -w -X main.version=$(V)" -o . github.com/klimby/version
	sudo chmod +x ./version

.PHONY: copy
copy: ## Copy from bin folder to root folder for use in this project
	cp ./bin/version ./

.PHONY: patch
patch: ## Patch version
	./version next --patch

.PHONY: minor
minor: ## Minor version
	./version next --minor

.PHONY: major
major: ## Major version
	./version next --major
