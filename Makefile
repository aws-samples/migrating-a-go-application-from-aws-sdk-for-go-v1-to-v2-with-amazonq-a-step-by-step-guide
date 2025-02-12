# Makefile for multiple Lambda functions (Windows and Mac compatible)

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
else
    DETECTED_OS := $(shell uname -s)
endif

# Variables
GOOS := linux
GOARCH := arm64
BINARY_NAME := bootstrap
LAMBDA_DIR := lambdafunction

# Lambda functions
LAMBDAS := main hitCounter updatePlayer

# Default target
all: $(LAMBDAS)

# Pattern rule for building and zipping Lambda functions
$(LAMBDAS):
ifeq ($(DETECTED_OS),Windows)
	@echo Building $(LAMBDA_DIR)\$@.go into $@.zip
	cd $(LAMBDA_DIR) && set GOOS=$(GOOS)& set GOARCH=$(GOARCH)& go build -o ..\$(BINARY_NAME) $@.go
	powershell Compress-Archive -Path $(BINARY_NAME) -DestinationPath $@.zip -Force
	del $(BINARY_NAME)
else
	@echo "Building $(LAMBDA_DIR)/$@.go into $@.zip"
	cd $(LAMBDA_DIR) && GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ../$(BINARY_NAME) $@.go
	zip $@.zip $(BINARY_NAME)
	rm $(BINARY_NAME)
endif

# Clean up
clean:
ifeq ($(DETECTED_OS),Windows)
	if exist $(BINARY_NAME) del $(BINARY_NAME)
	if exist *.zip del *.zip
else
	rm -f $(BINARY_NAME)
	rm -f *.zip
endif

# Get dependencies (if needed)
get-deps:
	go get github.com/aws/aws-sdk-go-v2
	go get github.com/aws/aws-sdk-go-v2/aws
	go get github.com/aws/aws-sdk-go-v2/config
	go get github.com/aws/aws-sdk-go-v2/service/dynamodb
	go get github.com/aws/aws-sdk-go-v2/service/dynamodb/types

# Phony targets
.PHONY: all clean get-deps $(LAMBDAS)
