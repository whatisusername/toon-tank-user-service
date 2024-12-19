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
	SHELL := /bin/bash
else ifeq ($(OS_NAME), Linux)
	SHELL := /bin/bash
else
	$(error Unsupported OS: $(OS_NAME))
endif

# Define variables
DOCKER_IMAGE_NAME := user-auth-handler
LAMBDA_NAME := userAuthHandler
TF_PATH := deployments/terraform

# Targets
.PHONY: build test local_run tf_init tf_plan tf_deploy docker_clean docker_prune gen_mock

build:
	docker build --platform linux/amd64 -t $(DOCKER_IMAGE_NAME):test -f ./deployments/docker/Dockerfile .

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

local_run:
	docker run --platform linux/amd64 -d --name $(DOCKER_IMAGE_NAME) -v ~/.aws-lambda-rie:/aws-lambda -p 9000:8080 \
		--entrypoint /aws-lambda/aws-lambda-rie $(DOCKER_IMAGE_NAME):test /main

tf_init:
ifeq ($(OS_NAME), Windows)
	@if exist "$(TF_PATH)\.terraform\" ( \
		rmdir /S /Q "$(TF_PATH)\.terraform" \
	)
	@( \
		echo env = "$(ENV)" \
	) > $(TF_PATH)\variables.auto.tfvars
else
	@if [[ -d "$(TF_PATH)/.terraform" ]]; then \
		rm -rf "$(TF_PATH)/.terraform"; \
	fi
	@cat $(TF_PATH)/variables.auto.tfvars <<EOF \
	env = "$(ENV)" \
	EOF
endif
	terraform -chdir=$(TF_PATH) init -backend-config='key=$(ENV)/lambda/$(LAMBDA_NAME)/terraform.tfstate'

tf_plan:
	terraform -chdir=$(TF_PATH) plan -out=tfplan

tf_deploy:
	terraform -chdir=$(TF_PATH) apply "tfplan"

docker_clean:
ifeq ($(OS_NAME), Windows)
	powershell -Command "docker ps -aq | ForEach-Object { docker stop $_ }; docker ps -aq | ForEach-Object { docker rm -f $_ }; docker images -aq | ForEach-Object { docker rmi -f $_ }"
else
	docker ps -aq | xargs -r docker stop && docker ps -aq | xargs -r docker rm -f && docker images -aq | xargs -r docker rmi -f
endif

docker_prune:
	docker system prune -a --volumes

gen_mock:
	mockery
