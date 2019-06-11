#!/bin/bash
set -Eeuo pipefail

# Templating
yaml2json test/fixtures/simple.yaml test/fixtures/simple.json
yaml2toml test/fixtures/simple.yaml test/fixtures/simple.toml
yaml2json test/fixtures/simple.part.yaml test/fixtures/simple.part.json
yaml2toml test/fixtures/simple.part.yaml test/fixtures/simple.part.toml
