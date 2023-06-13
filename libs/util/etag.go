package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"

	json "github.com/goccy/go-json"
	"github.com/pkg/errors"
)

// ETag generates an etag derived from a string source.
func ETag(input string) string {
	h := sha1.New()
	_, err := io.WriteString(h, input)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create new sha1 hash from string. input: '%s'", input))
	}

	byteArray := h.Sum(nil)

	encodedString := hex.EncodeToString(byteArray)

	return encodedString[:5]
}

// ETagAny generates an etag derived from any JSON marshalable source.
func ETagAny(input any) string {
	data, err := json.Marshal(input)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create new sha1 hash from input. input: '%v'", input))
	}

	return ETag(string(data))
}
