package validation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/util"
	"github.com/abtergo/abtergo/libs/validation"
)

func TestValidate(t *testing.T) {
	v := validation.NewValidator()

	validation.AddNotBeforeValidation(v)
	validation.AddEtagValidation(v)
	validation.AddPathValidation(v)

	type obj struct {
		Etag      string    `json:"etag" validate:"etag"`
		CreatedAt time.Time `json:"created_at,omitempty" validate:"not_before_date=2023-01-01"`
		Path      string    `json:"path,omitempty" validate:"path"`
	}

	tests := []struct {
		name    string
		obj     obj
		wantErr error
	}{
		{
			name: "valid",
			obj: obj{
				Etag:      "abc23",
				CreatedAt: util.MustParseDate("2023-01-01", time.DateOnly),
				Path:      "/path",
			},
			wantErr: nil,
		},
		{
			name: "empty etag",
			obj: obj{
				Etag:      "",
				CreatedAt: util.MustParseDate("2023-01-01", time.DateOnly),
			},
			wantErr: nil,
		},
		{
			name: "empty created at",
			obj: obj{
				Etag:      "",
				CreatedAt: util.MustParseDate("2023-01-01", time.DateOnly),
			},
			wantErr: nil,
		},
		{
			name: "invalid etag",
			obj: obj{
				Etag:      "ASBDJSI",
				CreatedAt: util.MustParseDate("2023-01-01", time.DateOnly),
			},
			wantErr: errors.New("Key: 'obj.etag' Error:Field validation for 'etag' failed on the 'etag' tag"),
		},
		{
			name: "created_at is before date",
			obj: obj{
				Etag:      "abc23",
				CreatedAt: util.MustParseDate("2021-01-01", time.DateOnly),
			},
			wantErr: errors.New("Key: 'obj.created_at' Error:Field validation for 'created_at' failed on the 'not_before_date' tag"),
		},
		{
			name: "invalid path",
			obj: obj{
				Etag:      "ASBDJSI",
				CreatedAt: util.MustParseDate("2023-01-01", time.DateOnly),
				Path:      "https://example.com/path",
			},
			wantErr: errors.New("Key: 'obj.etag' Error:Field validation for 'etag' failed on the 'etag' tag"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&tt.obj)
			if err == nil && tt.wantErr == nil {
				return
			}

			if err != nil && tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
				return
			}

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
