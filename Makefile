BINARY_NAME=marzban-api-gtw
BUILD_DIR=build

.PHONY: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go
