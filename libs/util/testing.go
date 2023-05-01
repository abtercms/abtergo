package util

import (
	"io"
	"net/http"
	"strings"
	"testing"

	json "github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

// DataToReaderHelper is a test helper function which creates a new io.Reader from any value that can be JSON marshalled.
func DataToReaderHelper(t *testing.T, v interface{}) io.Reader {
	data, err := json.Marshal(v)
	require.NoError(t, err)

	return strings.NewReader(string(data))
}

// ParseResponseHelper is a test helper function which can return any type from an HTTP response body.
func ParseResponseHelper[V any](t *testing.T, resp *http.Response, v V) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, v)
	require.NoError(t, err)
}
