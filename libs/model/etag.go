package model

import (
	"crypto/sha1"
	"encoding/hex"
	"io"

	json "github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func ETagFromString(input string) ETag {
	h := sha1.New()
	_, err := io.WriteString(h, input)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create new sha1 hash from string. input: '%s'", input))
	}

	byteArray := h.Sum(nil)

	encodedString := hex.EncodeToString(byteArray)

	return ETag(encodedString[:5])
}

// ETagFromAny generates an e-tag derived from any JSON marshalable source.
func ETagFromAny(input any) ETag {
	data, err := json.Marshal(input)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create new sha1 hash from input. input: '%v'", input))
	}

	return ETagFromString(string(data))
}
