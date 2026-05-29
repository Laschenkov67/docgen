package docx

import (
	"regexp"
	"strings"
)

func NormalizeRuns(src []byte) []byte {
	s := string(src)
	out := paragraphRe.ReplaceAllStringFunc(s, func(par string) string {
		if !needMerge(par) {
			return par
		}
		return mergeRuns(par)
	})
	return []byte(out)
}

var (
	paragraphRe = regexp.MustCompile(`(?s)<w:p\b[^>]*>.*?</w:p>`)
	wtRe        = regexp.MustCompile(`(?s)<w:t(\s[^>]*)?>(.*?)</w:t>`)
)

func needMerge(par string) bool {
	// Собираем чистый текст и смотрим, есть ли в нём незакрытые {{.
	var text strings.Builder
	for _, m := range wtRe.FindAllStringSubmatch(par, -1) {
		text.WriteString(m[2])
	}
	t := text.String()
	if strings.Count(t, "{{") == 0 {
		return false
	}
	for _, m := range wtRe.FindAllStringSubmatch(par, -1) {
		seg := m[2]
		o := strings.Count(seg, "{{")
		c := strings.Count(seg, "}}")
		if o != c {
			return true
		}
	}
	return false
}

func mergeRuns(par string) string {
	matches := wtRe.FindAllStringSubmatchIndex(par, -1)
	if len(matches) < 2 {
		return par
	}
	var combined strings.Builder
	for _, m := range matches {
		combined.WriteString(par[m[4]:m[5]])
	}

	var b strings.Builder
	last := 0
	for i, m := range matches {
		b.WriteString(par[last:m[0]])
		if i == 0 {
			b.WriteString(`<w:t xml:space="preserve">`)
			b.WriteString(combined.String())
			b.WriteString(`</w:t>`)
		} else {
			b.WriteString(`<w:t xml:space="preserve"></w:t>`)
		}
		last = m[1]
	}
	b.WriteString(par[last:])
	return b.String()
}
