.PHONY: tools proto test lint dependencies js-lib clean

CURRENT_DIR := $(shell pwd)

tools: ## Install Go tools
	@go install $(shell go list -f '{{join .Imports " "}}' tools/tools.go)

proto: ## Generate code from prot files
	buf generate

test: ## Run tests
	go test -v ./...

lint:
	tools/go-lint

dependencies: ## Install build dependencies
	$(CURRENT_DIR)/tools/download-deps.sh $(CURRENT_DIR)/tools/dependencies.toml

ui/lib/dist/index.js: ui/lib/node_modules
	cd ui/lib && yarn compile

ui/lib/dist/package.json:
	cp ui/lib/package.json ui/lib/dist

js-lib: ui/lib/dist/index.js ui/lib/dist/package.json

clean:
	rm -rf ui/lib/dist
