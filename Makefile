GO111MODULE=on

.PHONY: test
test:
	go test ./...

.PHONY: test.snap
test.snapshot:
	UPDATE_SNAPSHOTS=true go test ./...

.PHONY: lint
lint:
	./bin/golangci-lint run ./...

.PHONY: codegen
codegen:
	go generate ./...

.PHONY: fixturegen
fixturegen:
	./test/fixtures/generate-fixtures.sh
