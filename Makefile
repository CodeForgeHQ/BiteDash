.PHONY: lint fmt test

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

test:
	go test ./...