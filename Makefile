BINARY_NAME=github-actions-manager

.PHONY: build start clean lint

build:
	@go build -o bin/$(BINARY_NAME) cmd/manager/main.go
	@chmod +x bin/$(BINARY_NAME)

start: build
	@cd bin && ./$(BINARY_NAME) && cd ..

clean:
	@rm -rf bin/
	@go clean

lint:
	@golangci-lint run ./...
