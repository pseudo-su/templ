#!/bin/bash
set -Eeuo pipefail

## Templating
# Index
yaml2json --indent-json 2 --preserve-key-order test/fixtures/integrated/index.yaml test/fixtures/integrated/index.json
yaml2toml test/fixtures/integrated/index.yaml test/fixtures/integrated/index.toml
# Partials
yaml2json --indent-json 2 --preserve-key-order test/fixtures/integrated/partial.yaml test/fixtures/integrated/partial.json
yaml2toml test/fixtures/integrated/partial.yaml test/fixtures/integrated/partial.toml
