package fakeit

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/util"
)

func init() {
	AddDateRangeFaker()
	AddCSSURLFaker()
	AddJSURLFaker()
	AddPathFaker()
	AddWebsiteFaker()
	AddEtagFaker()
}

func AddDateRangeFaker() {
	gofakeit.AddFuncLookup("daterange2", gofakeit.Info{
		Category:    "abtergo",
		Description: "Date between start and end dates",
		Example:     "2023-01-01",
		Output:      "time.Time",
		Params: []gofakeit.Param{
			{Field: "start", Type: "string", Description: "Minimum date to use", Optional: true},
			{Field: "end", Type: "string", Description: "Maximum date to use", Optional: true},
			{Field: "layout", Type: "string", Description: "Date format to use", Optional: true, Default: time.DateOnly},
		},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			var startAt, endAt time.Time

			layout, err := info.GetString(m, "layout")
			if err != nil {
				return nil, errors.Wrap(err, "failed to get layout")
			}

			start, err := info.GetString(m, "start")
			if err == nil {
				startAt, err = time.Parse(layout, start)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to parse start date. layout: %s, end: %s", layout, start)
				}
			} else {
				startAt = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
			}

			end, err := info.GetString(m, "end")
			if err == nil {
				endAt, err = time.Parse(layout, end)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to parse end date. layout: %s, end: %s", layout, end)
				}
			} else {
				endAt = time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)
			}

			if endAt.Before(startAt) {
				return nil, fmt.Errorf("end date is before start date. start: %s, end: %s", startAt, endAt)
			}

			result := gofakeit.DateRange(startAt, endAt)

			if layout == time.DateOnly {
				return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, time.UTC), nil
			}

			return result, nil
		},
	})
}

func AddCSSURLFaker() {
	gofakeit.AddFuncLookup("url_css", gofakeit.Info{
		Category:    "abtergo",
		Description: "URL to a CSS file",
		Example:     "https://www.example.com/mypath.css",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			result := gofakeit.URL() + ".css"

			return result, nil
		},
	})
}

func AddJSURLFaker() {
	gofakeit.AddFuncLookup("url_js", gofakeit.Info{
		Category:    "abtergo",
		Description: "URL to a JS file",
		Example:     "https://www.example.com/mypath.js",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			result := gofakeit.URL() + ".js"

			return result, nil
		},
	})
}

func AddPathFaker() {
	gofakeit.AddFuncLookup("path", gofakeit.Info{
		Category:    "abtergo",
		Description: "Path to a file",
		Example:     "/mypath.js",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			// Slugs
			num := gofakeit.Number(1, 4)
			slug := make([]string, num)
			for i := 0; i < num; i++ {
				slug[i] = gofakeit.BS()
			}

			path := "/" + strings.Replace(strings.ToLower(strings.Join(slug, "/")), " ", "+", -1)

			if r.NormFloat64() < 0.5 {
				path += "." + gofakeit.FileExtension()
			}

			return path, nil
		},
	})
}

func AddWebsiteFaker() {
	gofakeit.AddFuncLookup("website", gofakeit.Info{
		Category:    "abtergo",
		Description: "Website URL",
		Example:     "https://sub.example.com/mypath-website.html",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			result := gofakeit.URL()

			return result, nil
		},
	})
}

func AddEtagFaker() {
	gofakeit.AddFuncLookup("etag", gofakeit.Info{
		Category:    "abtergo",
		Description: "E-tag",
		Example:     "aiso2",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			result := util.ETag(gofakeit.Word())

			return result, nil
		},
	})
}
