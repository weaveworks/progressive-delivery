.PHONY: proto test lint dependencies

CURRENT_DIR := $(shell pwd)

proto: ## Generate code from prot files
	buf generate

test: ## Run tests
	go test -v ./...

lint:
	tools/go-lint

dependencies: ## Install build dependencies
	$(CURRENT_DIR)/tools/download-deps.sh $(CURRENT_DIR)/tools/dependencies.toml
