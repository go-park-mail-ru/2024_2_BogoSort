BINARY_NAME=2024_2_BogoSort

BIN_DIR=bin

MAIN_PATH=./cmd/app

build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test ./...

test-cover:
	go test ./... -cover

clean:
	rm -rf $(BIN_DIR)

run: build
	./$(BIN_DIR)/$(BINARY_NAME)

.PHONY: build test clean run