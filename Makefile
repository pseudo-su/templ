GO111MODULE=on

.PHONY: default
default: gen lint test example.build

# TESTING

.PHONY: test
test:
	go test ./...

.PHONY: test.snap
test.snapshot:
	UPDATE_SNAPSHOTS=true go test ./...

# LINTING

.PHONY: lint
lint:
	./bin/golangci-lint run ./...

# CODE GENERATION

.PHONY: gen
gen: gen.code gen.fixtures

.PHONY: gen.code
gen.code:
	go generate ./...

.PHONY: gen.fixtures
gen.fixtures:
	./test/fixtures/generate-fixtures.sh

# EXAMPLE

.PHONY: example.build
example.build:
	go build -o stencil example/*.go
