package cleaner_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/abtergo/abtergo/libs/cleaner"
	mocks "github.com/abtergo/abtergo/mocks/libs/ablog"
)

func TestCleaner_Run(t *testing.T) {
	type fields struct {
		registry map[string]cleaner.Fn
	}
	type want struct {
		err error
		id  string
	}
	tests := []struct {
		name   string
		fields fields
		wants  []want
	}{
		{
			name: "empty",
		},
		{
			name: "no error",
			fields: fields{
				registry: map[string]cleaner.Fn{
					"foo": func() error {
						return nil
					},
				},
			},
		},
		{
			name: "simple error",
			fields: fields{
				registry: map[string]cleaner.Fn{
					"foo": func() error {
						return assert.AnError
					},
				},
			},
			wants: []want{
				{
					id:  "foo",
					err: assert.AnError,
				},
			},
		},
		{
			name: "multiple error",
			fields: fields{
				registry: map[string]cleaner.Fn{
					"foo": func() error {
						return assert.AnError
					},
					"bar": func() error {
						return assert.AnError
					},
					"baz": func() error {
						time.Sleep(10 * time.Millisecond)

						return assert.AnError
					},
					"quix": func() error {
						return nil
					},
				},
			},
			wants: []want{
				{
					id:  "foo",
					err: assert.AnError,
				},
				{
					id:  "bar",
					err: assert.AnError,
				},
				{
					id:  "baz",
					err: assert.AnError,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerMock := mocks.Logger{}

			for _, w := range tt.wants {
				loggerMock.EXPECT().Errorf(assert.AnError, mock.AnythingOfType("string"), w.id).Once()
			}

			c := cleaner.New(&loggerMock)

			for id, fn := range tt.fields.registry {
				c.Register(id, fn)
			}

			c.Run()

			loggerMock.AssertExpectations(t)
		})
	}
}

func TestCleaner_Unregister(t *testing.T) {
	t.Run("can remove a registered function", func(t *testing.T) {
		loggerMock := mocks.Logger{}

		c := cleaner.New(&loggerMock)

		c.Register("foo", func() error { panic("this must not be called") })

		c.Unregister("foo")

		c.Run()

		loggerMock.AssertExpectations(t)
	})
}
