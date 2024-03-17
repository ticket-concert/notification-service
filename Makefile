.PHONY: build clean development and deploy

install:
	@echo "Running download lib"
	go mod download

run:
	@echo "Running the application"
	go run cmd/main.go

unit-test:
	@echo "Running tests"
	mkdir -p ./test/coverage && \
		CGO_ENABLED=1 GOOS=linux go test $(BUILD_ARGS) -v ./... -coverprofile=./test/coverage/coverage.out

test-dev:
	@echo "Running tests"
	mkdir -p ./test/coverage && \
		CGO_ENABLED=1 go test -tags dynamic -v ./... -coverprofile=./test/coverage/coverage.out


coverage:
	@echo "Running tests with coverage"
	go tool cover -html=./test/coverage/coverage.out

build:
	@echo "Building the application"
	CGO_ENABLED=1 GOOS=linux go build $(BUILD_ARGS) -a -o build/bin/main cmd/main.go

start:
	@echo "Start the application"
	./build/bin/main

lint:
	@echo "Running linter"
	golangci-lint run

scan:
	@echo "Running security scan"
	gosec ./...

dev:
	go run -tags dynamic cmd/main.go
