package fakeit_test

import (
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/util"
)

func TestAddDateRangeFaker(t *testing.T) {
	fakeit.AddDateRangeFaker()

	t.Run("success", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01],[2023-12-31],[2006-01-02]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.GreaterOrEqual(t, f.Bar, util.MustParseDate("2023-01-01", time.DateOnly))
		assert.LessOrEqual(t, f.Bar, util.MustParseDate("2023-12-31", time.DateOnly))
	})

	t.Run("success with custom format", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023.01.01],[2023.12.31],[2006.01.02]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.GreaterOrEqual(t, f.Bar, util.MustParseDate("2023-01-01", time.DateOnly))
		assert.LessOrEqual(t, f.Bar, util.MustParseDate("2023-12-31", time.DateOnly))
	})

	t.Run("success with default format", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.GreaterOrEqual(t, f.Bar, util.MustParseDate("2023-01-01", time.DateOnly))
		assert.LessOrEqual(t, f.Bar, util.MustParseDate("2023-12-31", time.DateOnly))
	})

	t.Run("success with default end date and format", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.GreaterOrEqual(t, f.Bar, util.MustParseDate("2023-01-01", time.DateOnly))
	})

	t.Run("invalid start date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[foo]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed to parse start date."), err.Error())
	})

	t.Run("invalid end date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01],[foo]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed to parse end date."), err.Error())
	})

	t.Run("start date after end date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-02],[2023-01-01]}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("end date is before start date."), err.Error())
	})
}

func TestAddCSSURLFaker(t *testing.T) {
	fakeit.AddCSSURLFaker()

	validate := validator.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string `fake:"{url_css}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.NoError(t, validate.Var(f.Bar, "required,url"))
		assert.Equal(t, ".css", path.Ext(f.Bar))
	})
}

func TestAddJsURLFaker(t *testing.T) {
	fakeit.AddJSURLFaker()

	validate := validator.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string `fake:"{url_js}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.NoError(t, validate.Var(f.Bar, "required,url"))
		assert.Equal(t, ".js", path.Ext(f.Bar))
	})
}

func TestAddPathFaker(t *testing.T) {
	fakeit.AddPathFaker()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string `fake:"{path}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.Equal(t, path.Clean(f.Bar), f.Bar)
		assert.Regexp(t, regexp.MustCompile("^/[a-zA-Z0-9_/. +-]+$"), f.Bar)
	})
}

func TestAddWebsiteFaker(t *testing.T) {
	fakeit.AddWebsiteFaker()

	validate := validator.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string `fake:"{website}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.NoError(t, validate.Var(f.Bar, "required,url"))
	})
}
