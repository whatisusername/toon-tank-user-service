# Detect the operating system
ifeq ($(OS), Windows_NT)
	OS_NAME := Windows
else
	OS_NAME := $(shell uname -s)
endif

# Include the appropriate Makefile based on the OS
ifeq ($(OS_NAME), Windows)
	SHELL := cmd
else ifeq ($(OS_NAME), Darwin)
	SHELL := /bin/sh
else ifeq ($(OS_NAME), Linux)
	SHELL := /bin/bash
else
	$(error Unsupported OS: $(OS_NAME))
endif

# Targets
.PHONY: test gen_mock

test:
ifeq ($(OS_NAME), Windows)
	@if not exist "docs\" ( \
		mkdir docs \
	)
else
	@if [[ -d "docs" ]]; then \
		mkdir docs; \
	fi
endif
	go test ./... -v -coverprofile="docs/coverage.out"
	go tool cover -html="docs/coverage.out" -o "docs/coverage.html"

gen_mock:
	mockery
