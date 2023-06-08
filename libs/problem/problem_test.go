package problem

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/arr"
)

func TestFromError(t *testing.T) {
	type args struct {
		baseURL string
		err     error
	}
	tests := []struct {
		name string
		args args
		want Problem
	}{
		{
			name: "an error",
			args: args{
				baseURL: "https://example.com/",
				err:     assert.AnError,
			},
			want: Problem{
				Type:   "https://example.com/problem/unknown-error",
				Title:  arr.UnknownError.GetTitle(),
				Detail: "assert.AnError general error for testing",
				Status: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid user input",
			args: args{
				baseURL: "https://example.com",
				err:     arr.Wrap(arr.InvalidUserInput, assert.AnError, "foo"),
			},
			want: Problem{
				Type:   "https://example.com/problem/invalid-user-input",
				Title:  arr.InvalidUserInput.GetTitle(),
				Detail: "foo: assert.AnError general error for testing.",
				Status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FromError(tt.args.baseURL, tt.args.err))
		})
	}
}
