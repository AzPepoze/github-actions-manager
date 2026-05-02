BINARY_NAME=actions-manager

.PHONY: build start clean

build:
	@go build -o bin/$(BINARY_NAME) cmd/manager/main.go
	@chmod +x bin/$(BINARY_NAME)
	@echo "Build success!!"

start: build
	@./bin/$(BINARY_NAME)

clean:
	@rm -rf bin/
	@go clean
	@echo "Clean success!!"
