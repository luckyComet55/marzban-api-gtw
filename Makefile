BINARY_NAME=marzban-api-gtw
BUILD_DIR=build
GEN_DIR=contract

.PHONY: build generate clean help

generate:
	protoc -I proto proto/contract.proto --go_out=./$(GEN_DIR)/ \
	--go_opt=paths=source_relative --go-grpc_out=./$(GEN_DIR) \
	--go-grpc_opt=paths=source_relative

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go

setup:
	go mod download

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(GEN_DIR)

help:
	@echo "Available commands:"
	@echo "    build        - builds Marzban client binary"
	@echo "    clean        - cleans build artifacts"
