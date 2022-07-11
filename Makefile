.PHONY: proto test lint tools dependencies js-lib clean dev-server

CURRENT_DIR := $(shell pwd)

proto: ## Generate code from prot files
	buf generate

test: ## Run tests
	go test -v ./...

lint:
	$(CURRENT_DIR)/tools/go-lint

tools: ## Install Go tools
	go install $(shell go list -f '{{join .Imports " "}}' tools/tools.go)

dependencies: tools ## Install build dependencies
	$(CURRENT_DIR)/tools/download-deps.sh $(CURRENT_DIR)/tools/dependencies.toml

create-dev-cluster:
	$(CURRENT_DIR)/tools/bin/kind create cluster --name pd-dev

delete-dev-cluster:
	$(CURRENT_DIR)/tools/bin/kind delete cluster --name pd-dev

dev-cluster: ## Start a dev server with Tilt
	$(CURRENT_DIR)/tools/bin/tilt up

ui/lib/dist/index.js: ui/lib/node_modules
	cd ui/lib && yarn compile

ui/lib/dist/package.json: ui/lib/package.json
	cp ui/lib/package.json ui/lib/dist

js-lib: ui/lib/dist/index.js ui/lib/dist/package.json

clean:
	rm -rf ui/lib/dist
