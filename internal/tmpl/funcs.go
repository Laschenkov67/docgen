package tmpl

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/laschenkov67/docgen/format"
)

func DefaultFuncs() template.FuncMap {
	return template.FuncMap{
		"rubles": format.Rubles,
		"money":  format.Money,
		"sum": func(items any, field string) float64 {
			return format.SumField(items, field)
		},

		"date":     format.Date,
		"dateLong": format.DateLong,
		"now":      func() time.Time { return time.Now() },

		"upper":   strings.ToUpper,
		"lower":   strings.ToLower,
		"title":   strings.Title,
		"trim":    strings.TrimSpace,
		"replace": strings.ReplaceAll,
		"default": func(def, v any) any {
			if v == nil || fmt.Sprintf("%v", v) == "" {
				return def
			}
			return v
		},

		"inn": format.FormatINN,
		"kpp": format.FormatKPP,

		"inc": func(i int) int { return i + 1 },
		"add": func(a, b float64) float64 { return a + b },
		"mul": func(a, b float64) float64 { return a * b },
	}
}
