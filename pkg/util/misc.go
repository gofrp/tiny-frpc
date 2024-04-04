package util

import (
	"bytes"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode"
)

func JSONEncode(args interface{}) string {
	buf, err := json.Marshal(args)
	if err != nil {
		return fmt.Sprintf("args: %v, error: %v", args, err)
	}
	return string(buf)
}

func GetAuthKey(token string, timestamp int64) (key string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(token))
	md5Ctx.Write([]byte(strconv.FormatInt(timestamp, 10)))
	data := md5Ctx.Sum(nil)
	return hex.EncodeToString(data)
}

func ConstantTimeEqString(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
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

func handleTermSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
