package tmpl

import (
	"fmt"
	"strings"
	"text/template"
	"time"
	"unicode"

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
		"title":   titleCase,
		"trim":    strings.TrimSpace,
		"replace": strings.ReplaceAll,
		"default": func(def, v any) any {
			if v == nil || fmt.Sprintf("%v", v) == "" {
				return def
			}
			return v
		},

		"inn": format.INN,
		"kpp": format.KPP,

		"inc": func(i int) int { return i + 1 },
		"add": func(a, b float64) float64 { return a + b },
		"mul": func(a, b float64) float64 { return a * b },
	}
}

// titleCase — замена устаревшей strings.Title с тем же поведением (заглавная
// буква после любого не буквенно-цифрового разделителя), но без deprecated API.
func titleCase(s string) string {
	prevIsLetterOrDigit := false
	return strings.Map(func(r rune) rune {
		if !prevIsLetterOrDigit {
			prevIsLetterOrDigit = unicode.IsLetter(r) || unicode.IsDigit(r)
			return unicode.ToTitle(r)
		}
		prevIsLetterOrDigit = unicode.IsLetter(r) || unicode.IsDigit(r)
		return r
	}, s)
}
