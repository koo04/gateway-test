SHELL := /bin/bash

TARGET := gateway-test
.DEFAULT_GOAL: $(TARGET)

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean fmt lint test test-all serve

all: clean build

start: serve

serve: build
	@./bin/$(TARGET)

$(TARGET): $(SRC)
	go build $(LDFLAGS) -o bin/$(TARGET) -v main.go

build: $(TARGET)
	@true

clean:
	@rm -f ./bin/*
	@rm -rf ./internal/gen/proto/go/*

pb: pb-grpc-gateway
	@true

pb-grpc-gateway:
	# building proto go grpc gateway modules...
	@docker run --rm \
        --platform linux/amd64 \
        -v ${PWD}/./:/app \
        -e INPUT_PATH=./testapis/koo04 \
        -e OUTPUT_PATH=./ \
        -e PROTO_FILE=api/v1/test.proto \
        -e GEN_GO=true \
        -e GEN_GO_GRPC=true \
        -e GEN_GRPC_GATEWAY=true \
        -e GRPC_API_CONFIGURATION=./testapis/koo04/api/v1/test.yaml \
        koothegreat/proto-builder:latest
