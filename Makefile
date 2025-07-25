BINARY_NAME=marzban-api-gtw
BUILD_DIR=build

.PHONY: build clean help

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go

setup:
	go mod download

clean:
	rm -rf $(BUILD_DIR)

help:
	@echo "Available commands:"
	@echo "    build        - builds Marzban client binary"
	@echo "    clean        - cleans build artifacts"
