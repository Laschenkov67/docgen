package format

import "strings"

func ValidINN(s string) bool {
	s = strings.TrimSpace(s)
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	switch len(s) {
	case 10:
		return innCheck(s, []int{2, 4, 10, 3, 5, 9, 4, 6, 8, 0}, 9)
	case 12:
		ok1 := innCheck(s, []int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0, 0}, 10)
		ok2 := innCheck(s, []int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0}, 11)
		return ok1 && ok2
	}
	return false
}

func innCheck(s string, w []int, pos int) bool {
	sum := 0
	for i := 0; i < len(w); i++ {
		sum += int(s[i]-'0') * w[i]
	}
	return sum%11%10 == int(s[pos]-'0')
}

func ValidKPP(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) != 9 {
		return false
	}
	for i, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
		case i >= 4 && i <= 5 && ((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')):
		default:
			return false
		}
	}
	return true
}

func FormatINN(s string) string { return strings.TrimSpace(s) }

func FormatKPP(s string) string { return strings.TrimSpace(s) }
