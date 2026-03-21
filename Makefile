BIN_DIR := bin
BINARY := agdev

.PHONY: build install run-version release-patch release-minor release-major clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

install:
	go install .

run-version:
	go run . version

release-patch:
	./scripts/release.sh patch

release-minor:
	./scripts/release.sh minor

release-major:
	./scripts/release.sh major

clean:
	rm -rf $(BIN_DIR)
