BIN_NAME := tracer

BIN_DIR := bin

SRC_DIR := cmd/app/main.go

GO_BUILD := go build -o $(BIN_DIR)/$(BIN_NAME) $(SRC_DIR)

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build  - Build the application"
	@echo "  clean  - Remove the built binary"
	@echo "  run    - Build and run the application"
	@echo "  help   - Show this help message"

run: build
	./$(BIN_DIR)/$(BIN_NAME)

build:
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD)

clean:
	rm -f $(BIN_DIR)/$(BIN_NAME)

