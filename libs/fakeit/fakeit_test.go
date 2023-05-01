package fakeit_test

import (
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/fakeit"
)

func TestAddDateRangeFaker(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.GreaterOrEqual(t, f.Bar, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	})

	t.Run("missing start date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed retreiving the start param."), err.Error())
	})

	t.Run("invalid start date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[foo]}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed parsing the start date."), err.Error())
	})

	t.Run("missing end date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01]}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed retreiving the end param."), err.Error())
	})

	t.Run("invalid end date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-01],[foo]}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("failed parsing the end date."), err.Error())
	})

	t.Run("start date after end date", func(t *testing.T) {
		type foo struct {
			Bar time.Time `fake:"{daterange2:[2023-01-02],[2023-01-01]}"`
		}

		fakeit.AddDateRangeFaker()

		f := foo{}

		err := gofakeit.Struct(&f)
		require.Error(t, err)

		assert.Regexp(t, regexp.MustCompile("start date after end date."), err.Error())
	})
}

func TestAddCssURLFaker(t *testing.T) {
	type foo struct {
		Bar string `fake:"{url_css}"`
	}

	fakeit.AddCSSURLFaker()

	f := foo{}

	err := gofakeit.Struct(&f)
	require.NoError(t, err)

	assert.NotEmpty(t, f.Bar)
	assert.Regexp(t, regexp.MustCompile("^https?://[a-zA-Z0-9_/.-]+.css$"), f.Bar)
}

func TestAddJsURLFaker(t *testing.T) {
	type foo struct {
		Bar string `fake:"{url_js}"`
	}

	fakeit.AddJsURLFaker()

	f := foo{}

	err := gofakeit.Struct(&f)
	require.NoError(t, err)

	assert.NotEmpty(t, f.Bar)
	assert.Regexp(t, regexp.MustCompile("^https?://[a-zA-Z0-9_/.-]+.js$"), f.Bar)
}

func TestAddPathFaker(t *testing.T) {
	type foo struct {
		Bar string `fake:"{path}"`
	}

	fakeit.AddPathFaker()

	f := foo{}

	err := gofakeit.Struct(&f)
	require.NoError(t, err)

	assert.NotEmpty(t, f.Bar)
	assert.Equal(t, path.Clean(f.Bar), f.Bar)
	assert.Regexp(t, regexp.MustCompile("^/[a-zA-Z0-9_/-]+$"), f.Bar)
}

func TestAddWebsiteFaker(t *testing.T) {
	type foo struct {
		Bar string `fake:"{website}"`
	}

	fakeit.AddWebsiteFaker()

	f := foo{}

	err := gofakeit.Struct(&f)
	require.NoError(t, err)

	assert.NotEmpty(t, f.Bar)
	assert.Regexp(t, regexp.MustCompile("^https://[a-zA-Z0-9_/.-]+$"), f.Bar)
}

func TestAddEtagFaker(t *testing.T) {
	type foo struct {
		Bar string `fake:"{etag}"`
	}

	fakeit.AddEtagFaker()

	f := foo{}

	err := gofakeit.Struct(&f)
	require.NoError(t, err)

	assert.NotEmpty(t, f.Bar)
	const pattern = "^[a-z0-9]{5}$"
	assert.Regexp(t, regexp.MustCompile(pattern), f.Bar)
}
