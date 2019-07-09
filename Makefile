GO111MODULE=on

.PHONY: default
default: gen lint test

.PHONY: test
test:
	go test ./...

.PHONY: test.snap
test.snapshot:
	UPDATE_SNAPSHOTS=true go test ./...

.PHONY: lint
lint:
	./bin/golangci-lint run ./...

.PHONY: gen
gen: gen.code gen.fixtures

.PHONY: gen.code
gen.code:
	go generate ./...

.PHONY: gen.fixtures
gen.fixtures:
	./test/fixtures/generate-fixtures.sh
