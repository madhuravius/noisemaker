deps:
	go mod download
.PHONY: deps

init: deps
	go generate ./...
.PHONY: init

build-image:
	podman build -t noisemaker .
.PHONY: build-image

build:
	go build -ldflags="-s -w" -o ./bin/noisemaker main.go
.PHONY: build

lint:
	podman run \
		--rm \
		-v $(shell pwd):/app \
		-w /app \
		docker.io/golangci/golangci-lint:v1.52 \
		golangci-lint run
.PHONY: lint

pretty:
	go fmt ./...
.PHONY: pretty
