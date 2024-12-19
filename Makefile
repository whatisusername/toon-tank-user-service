# Detect the operating system
ifeq ($(OS), Windows_NT)
	OS_NAME := Windows
else
	OS_NAME := $(shell uname -s)
endif

# Include the appropriate Makefile based on the OS
ifeq ($(OS_NAME), Windows)
	SHELL := cmd
	include Makefile.windows
else ifeq ($(OS_NAME), Darwin)
	SHELL := /bin/sh
	include Makefile.mac
else ifeq ($(OS_NAME), Linux)
	SHELL := /bin/bash
	include Makefile.linux
else
	$(error Unsupported OS: $(OS_NAME))
endif

.PHONY: build run clean
APP_EXECUTABLE := main
build:
	go build -o $(APP_EXECUTABLE) cmd/main.go

run: build
	./$(APP_EXECUTABLE)

clean:
	go clean
	rm $(APP_EXECUTABLE)
