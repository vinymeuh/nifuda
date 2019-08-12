SHELL := $(shell which bash)
ENV = /usr/bin/env

.SHELLFLAGS = -c

.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:

.PHONY: test
.DEFAULT_GOAL := help

VERSION = `git describe --tags --always`
BUILD   = `date +%FT%T%z`

LDFLAGS = -w -s -X main.version=${VERSION} -X main.build=${BUILD}

build: clean ## Build binary
	#go build -ldflags "${LDFLAGS}"
	@for cmd in $$(ls cmd); do \
		if [ -f cmd/$$cmd/main.go ]; then \
			echo "building $$cmd"; \
			go build -ldflags "${LDFLAGS}" github.com/vinymeuh/nifuda-ng/cmd/$$cmd; \
		fi; \
	done;

coverage: ## Create coverage report
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt

clean: ## Delete binary
	rm -f nifuda

help: ## Show Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

install: clean ## Install binary in GOPATH
	go install -ldflags "${LDFLAGS}"

test: ## Run tests
	go test -coverprofile=coverage.txt -ldflags "${LDFLAGS}" ./pkg/*
