#!/bin/sh
set -eu

cd "$(dirname "$(realpath "$0")")"

DEST="${1:?usage: build-tools.sh <dest>}"
# Anchor relative dests to the project root: builds run inside src/<tool>,
# where a relative -o would otherwise land beside the module instead of in dest.
case "$DEST" in
/*) ;;
*) DEST="$(pwd)/$DEST" ;;
esac

mkdir -p "$DEST/manifests"

for tool_dir in src/*/; do
	name="$(basename "$tool_dir")"
	[ -f "$tool_dir/go.mod" ] || continue
	# src/shared is the shared runtime lib, not a shipped tool.
	[ "$name" = "shared" ] && continue

	(
		cd "$tool_dir"
		go build -o "$DEST/$name" .
	)

	manifest="manifests/$name.yaml"
	[ -f "$manifest" ] || {
		echo "build-tools.sh: $tool_dir has go.mod but $manifest is missing" >&2
		exit 1
	}
	install -m 0644 "$manifest" "$DEST/$manifest"
done
