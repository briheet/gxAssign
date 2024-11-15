test:
	@go test -v -race -timeout 1s ./...

build:
	@go build ./cmd/gx

lint:
	@golangci-lint run

.PHONY: test build lint
