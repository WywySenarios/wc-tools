# Contributing to wc-tools

## Tool conventions

Every shipped tool:

- is a Go module (`src/<tool>/go.mod`).
- produces a static binary.
- is a JSON tool: args arrive as one JSON object in a single argv element, or on stdin as
  `<tool> --json`
- reads tool config from `$WC_CONFIG_DIR/<tool>/config.yaml` — a single env var, no
  fallback chain. Missing config means built-in defaults — never an error.

Tools share an internal library for parsing input.

## Layout

- `src/shared/*`
- `src/<tool>/*`
- `config/schema/tools.yaml`
- `build-tools.sh`, deployment contract (builds each Go module with `go build` and stages the tool tree into `<dest>`)
- `scripts/tests/`, test suite scripts
- `manifests/<tool>.yaml`

## Pre-commit

This repo uses pre-commit. Please run `pre-commit install` to install the git hooks.
