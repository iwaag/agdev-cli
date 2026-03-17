BIN_DIR := bin
BINARY := agdev

.PHONY: build install run-version clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

install:
	go install .

run-version:
	go run . version

clean:
	rm -rf $(BIN_DIR)
