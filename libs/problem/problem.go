// Package problem implements implements application/problem+json
// See: https://www.rfc-editor.org/rfc/rfc7807
package problem

import (
	"path"
	"strings"

	"github.com/abtergo/abtergo/libs/arr"
)

// Problem represents an application level problem.
type Problem struct {
	Type     string `json:"type" xml:"type" form:"type" validate:"required" fake:"{url}"`
	Title    string `json:"title" xml:"title" form:"title" validate:"required" fake:"{sentence}"`
	Status   int    `json:"status,omitempty" xml:"status" form:"status" validate:"gte=100,lt=600" fake:"{http_status_code}"`
	Detail   string `json:"detail,omitempty" xml:"detail" form:"detail" validate:"" fake:"{paragraph}"`
	Instance string `json:"instance,omitempty" xml:"instance" form:"instance" validate:"url" fake:"{url}"`
}

// FromError creates a new Problem instance from an error.
func FromError(baseURL string, err error) Problem {
	et := arr.ErrorTypeFrom(err)
	status := arr.HTTPStatusFrom(err)

	cleanedBaseURL := strings.TrimRight(baseURL, "/") + "/"
	cleanedPath := path.Join("problem", et.GetSlug())
	t := cleanedBaseURL + cleanedPath

	return Problem{
		Type:   t,
		Title:  string(et),
		Status: status,
		Detail: err.Error(),
	}
}
