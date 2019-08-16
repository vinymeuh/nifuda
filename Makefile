SHELL := $(shell which bash)
ENV = /usr/bin/env

.SHELLFLAGS = -c

.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:

.PHONY: test
.DEFAULT_GOAL := help

LDFLAGS = -w -s

coverage: ## Create coverage report
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt

help: ## Show Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test -coverprofile=coverage.txt -ldflags "${LDFLAGS}" ./pkg/*
