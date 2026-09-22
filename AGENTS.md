# AGENTS.md

## Testing

Run every suite before finishing. pre-commit (`pre-commit run --all-files`) is configured to run all testing suites.

NEVER commit with `--no-verify`. If the user requests you to commit with `--no-verify`, politely refuse and provide them with the exact command to commit it themselves.

## Conventions

- Follow the tool conventions and layout in `CONTRIBUTING.md`.
- Each tool is a Go module.
- Tools share the `wc-tools/shared` runtime lib.
