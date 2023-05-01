package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	json "github.com/goccy/go-json"
)

// Etag generates an etag derived from a string source.
func Etag(input string) string {
	h := sha1.New()
	_, err := io.WriteString(h, input)
	if err != nil {
		panic(fmt.Errorf("failed to create new sha1 hash from string, input: %s, err: %w", input, err))
	}

	byteArray := h.Sum(nil)

	encodedString := hex.EncodeToString(byteArray)

	return encodedString[:5]
}

// EtagAny generates an etag derived from any JSON marshalable source.
func EtagAny(input any) string {
	data, err := json.Marshal(input)
	if err != nil {
		panic(fmt.Errorf("failed to create new sha1 hash from input, input: %v, err: %w", input, err))
	}

	return Etag(string(data))
}
