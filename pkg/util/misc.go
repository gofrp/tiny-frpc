package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"
)

func JSONEncode(args interface{}) string {
	buf, err := json.Marshal(args)
	if err != nil {
		return fmt.Sprintf("args: %v, error: %v", args, err)
	}
	return string(buf)
}

func EmptyOr[T comparable](v T, fallback T) T {
	var zero T
	if zero == v {
		return fallback
	}
	return v
}

// IsJSONBuffer scans the provided buffer, looking
// for an open brace indicating this is JSON.
func IsJSONBuffer(buf []byte) bool {
	return hasJSONPrefix(buf)
}

var jsonPrefix = []byte("{")

// hasJSONPrefix returns true if the provided buffer appears to start with
// a JSON open brace.
func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is
// prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

func Ternary[T any](condition bool, ifOutput T, elseOutput T) T {
	if condition {
		return ifOutput
	}

	return elseOutput
}
