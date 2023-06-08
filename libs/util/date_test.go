package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/util"
)

func TestCloneDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := time.Now()

		actual := util.CloneDate(expected)

		assert.Equal(t, expected.UnixNano(), actual.UnixNano())
		assert.NotSame(t, expected, actual)
	})
}

func TestMustParseDate(t *testing.T) {
	t.Run("failure in case of layout not matching value", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			util.MustParseDate("2016-02-23", time.RFC822)
		})
	})

	t.Run("failure in case of invalid layout", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			util.MustParseDate("2016-02-23", "foo")
		})
	})

	t.Run("failure in case of invalid date", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			util.MustParseDate("foor", time.DateOnly)
		})
	})

	type args struct {
		date   string
		format string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "success",
			args: args{
				date:   "2020-01-01",
				format: "2006-01-02",
			},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, util.MustParseDate(tt.args.date, tt.args.format), "MustParseDate(%v, %v)", tt.args.date, tt.args.format)
		})
	}
}
