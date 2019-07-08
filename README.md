# Templ

Templ is a library to allow loading structured data (YAML/JSON/TOML) and applying go-like templating to the string fields, enabling relative references to the same document.

EXAMPLE INPUT:

```yaml
data:
  string: Hello World
  array:
    - Array
    - of
    - strings
  object:
    key: val
example_1: ${ self `data.string` }
example_2: ${ self `data.array` }
example_2: ${ self `data.object` }
```

OUTPUT:

```yaml
data:
  string: Hello World
  array:
    - Array
    - of
    - strings
  object:
    key: val
example_1: Hello World
example_2:
  array:
    - Array
    - of
    - strings
example_2:
  object:
    key: val
```

## Contributing

### Install Tool Dependencies

```shell
# install golanglint-ci into ./bin
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.16.0
# Install gobin globally
env GO111MODULE=off go get -u github.com/myitcv/gobin
```

### Build tasks

**NOTE**: if this repo is cloned into you `GOPATH` you will need to prefix all commands with `GO111MODULES=on`

* test: `go test ./...`
* lint: `./bin/golangci-lint run ./...`
* codegen: `go generate ./...` [link](https://github.com/go-swagger/go-swagger/issues/1724#issuecomment-469335593)
