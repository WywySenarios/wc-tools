// Package shared provides the input contract shared by every wc-tools tool.
package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// MaxArgsBytes bounds a stdin-supplied argument object.
const MaxArgsBytes = 1 << 20

// ReadArgs resolves the JSON argument object for a tool invocation and decodes
// it into v.
//
// args is os.Args[1:]. The object arrives either as the single remaining
// argument, or on stdin when invoked as `<tool> --json`. Tools accept exactly
// one argument source; every other shape is an error.
func ReadArgs(args []string, stdin io.Reader, v any) error {
	var raw []byte
	switch {
	case len(args) == 1 && args[0] == "--json":
		data, err := io.ReadAll(io.LimitReader(stdin, MaxArgsBytes+1))
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		if len(data) > MaxArgsBytes {
			return fmt.Errorf("stdin exceeds %d bytes", MaxArgsBytes)
		}
		raw = data
	case len(args) == 1:
		raw = []byte(args[0])
	case len(args) == 0:
		return errors.New("missing arguments: pass one JSON object or --json")
	default:
		return errors.New("expected exactly one JSON argument or --json")
	}

	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return errors.New("empty argument object")
	}
	if raw[0] != '{' {
		return errors.New("arguments must be a JSON object")
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode arguments: %w", err)
	}

	var extra any
	switch err := dec.Decode(&extra); {
	case err == io.EOF:
	case err == nil:
		return errors.New("arguments contain trailing data")
	default:
		return fmt.Errorf("arguments contain trailing data: %w", err)
	}

	return nil
}
