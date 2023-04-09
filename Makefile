init:
	go mod download
	go generate ./...
.PHONY: init

build-image:
	podman build -t trashbin .
.PHONY: build-image

build:
	go build -ldflags="-s -w" -o bin/trashbin main.go
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
