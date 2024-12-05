# Makefile for multiple Lambda functions

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
	@echo "Building $(LAMBDA_DIR)/$@.go into $@.zip"
	cd $(LAMBDA_DIR) && GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ../$(BINARY_NAME) $@.go
	zip $@.zip $(BINARY_NAME)
	rm $(BINARY_NAME)

# Clean up
clean:
	rm -f $(BINARY_NAME)
	rm -f *.zip

# Get dependencies (if needed)
get-deps:
	go get github.com/aws/aws-sdk-go-v2
	go get github.com/aws/aws-sdk-go-v2/aws
	go get github.com/aws/aws-sdk-go-v2/config
	go get github.com/aws/aws-sdk-go-v2/service/dynamodb
	go get github.com/aws/aws-sdk-go-v2/service/dynamodb/types

# Phony targets
.PHONY: all clean get-deps $(LAMBDAS)
