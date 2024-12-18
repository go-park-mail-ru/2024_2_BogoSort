BINARY_NAME=2024_2_BogoSort

BIN_DIR=bin

MAIN_PATH=./cmd/app

GOTMPDIR=./tmp

EASYJSON=easyjson

build:
	mkdir -p $(BIN_DIR)
	GOTMPDIR=$(GOTMPDIR) go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test:
	GOTMPDIR=$(GOTMPDIR) go test ./...

test-cover:
	GOTMPDIR=$(GOTMPDIR) go test ./... -cover

clean:
	rm -rf $(BIN_DIR) $(GOTMPDIR)

run: build
	./$(BIN_DIR)/$(BINARY_NAME)

swagger:
	swag init -g $(MAIN_PATH)/main.go -o ./docs

generate:
	@echo "Generating easyjson files..."
	@find . -name '*.go' -exec $(EASYJSON) -all {} \;

.PHONY: build test clean run swagger