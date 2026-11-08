# Go command.
GO ?= go

# Tool versions.
GCL_VERSION  := v2.15.0
SWAG_VERSION := v1.16.4

# Project paths.
CMD_PATH    := ./cmd/shimakaze
BINARY_NAME := shimakaze
BINARY_PATH := $(CMD_PATH)/$(BINARY_NAME)

.PHONY: fmt clean lint swagger install-linter install-swagger \
	build api consumer cron-update cron-fill

##@ Tools

fmt: ## Standarize go coding style for the whole project.
	@$(GO) fmt ./...

clean: ## Clean project binary, test, and coverage file.
	@$(GO) clean ./...
	@rm -f $(COVER_OUT)

lint: fmt ## Lint go source code.
	@golangci-lint run -D errcheck --timeout 5m

swagger: ## Generate swagger docs.
	@swag init -g $(CMD_PATH)/main.go -o ./docs

install-linter: ## Install golangci-lint.
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GCL_VERSION)

install-swagger: ## Install swag.
	@$(GO) install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

##@ Build & Run

build: clean fmt ## Build the project executable binary. 
	@$(GO) build -o $(BINARY_PATH) -v $(CMD_PATH)

api: build ## Build and run API server.
	@$(BINARY_PATH) server

consumer: build ## Build and run consumer.
	@$(BINARY_PATH) consumer

cron-update: build ## Build and run cron update.
	@$(BINARY_PATH) cron update

cron-fill: build ## Build and run cron fill.
	@$(BINARY_PATH) cron fill