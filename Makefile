.PHONY: default
default:
	make build

.PHONY: build ## Builds and ensures all of the source files in this repo are valid.
build:
	go build -o bin/yamlsplit

.PHONY: test ## Runs unit tests for this repo.
test:
	go test ./...
