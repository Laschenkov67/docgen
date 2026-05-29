package format

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

func Money(v float64) string {
	neg := v < 0
	if neg {
		v = -v
	}
	intPart := int64(math.Floor(v))
	frac := int64(math.Round((v - float64(intPart)) * 100))
	if frac == 100 {
		intPart++
		frac = 0
	}
	s := fmt.Sprintf("%d", intPart)
	// группировка по 3
	var b strings.Builder
	n := len(s)
	for i, ch := range s {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(ch)
	}
	res := fmt.Sprintf("%s,%02d", b.String(), frac)
	if neg {
		return "-" + res
	}
	return res
}

func Rubles(v float64) string {
	if v < 0 {
		return "минус " + Rubles(-v)
	}
	rub := int64(math.Floor(v))
	kop := int64(math.Round((v - float64(rub)) * 100))
	if kop == 100 {
		rub++
		kop = 0
	}
	words := intToWords(rub, genderMale)
	rubWord := plural(rub, "рубль", "рубля", "рублей")
	kopWord := plural(kop, "копейка", "копейки", "копеек")
	return fmt.Sprintf("%s %s %02d %s",
		capitalize(words), rubWord, kop, kopWord)
}

const (
	genderMale   = 0
	genderFemale = 1
)

var (
	units0_19m = []string{"ноль", "один", "два", "три", "четыре", "пять", "шесть",
		"семь", "восемь", "девять", "десять", "одиннадцать", "двенадцать",
		"тринадцать", "четырнадцать", "пятнадцать", "шестнадцать",
		"семнадцать", "восемнадцать", "девятнадцать"}
	units0_19f = []string{"ноль", "одна", "две", "три", "четыре", "пять", "шесть",
		"семь", "восемь", "девять", "десять", "одиннадцать", "двенадцать",
		"тринадцать", "четырнадцать", "пятнадцать", "шестнадцать",
		"семнадцать", "восемнадцать", "девятнадцать"}
	tens = []string{"", "", "двадцать", "тридцать", "сорок", "пятьдесят",
		"шестьдесят", "семьдесят", "восемьдесят", "девяносто"}
	hundreds = []string{"", "сто", "двести", "триста", "четыреста", "пятьсот",
		"шестьсот", "семьсот", "восемьсот", "девятьсот"}
)

func tripletToWords(n int64, gender int) string {
	if n == 0 {
		return ""
	}
	var parts []string
	h := n / 100
	rem := n % 100
	if h > 0 {
		parts = append(parts, hundreds[h])
	}
	if rem < 20 {
		if rem > 0 {
			if gender == genderFemale {
				parts = append(parts, units0_19f[rem])
			} else {
				parts = append(parts, units0_19m[rem])
			}
		}
	} else {
		t := rem / 10
		u := rem % 10
		parts = append(parts, tens[t])
		if u > 0 {
			if gender == genderFemale {
				parts = append(parts, units0_19f[u])
			} else {
				parts = append(parts, units0_19m[u])
			}
		}
	}
	return strings.Join(parts, " ")
}

func intToWords(n int64, gender int) string {
	if n == 0 {
		return "ноль"
	}
	scales := []struct {
		div    int64
		gender int
		one    string
		few    string
		many   string
	}{
		{1_000_000_000, genderMale, "миллиард", "миллиарда", "миллиардов"},
		{1_000_000, genderMale, "миллион", "миллиона", "миллионов"},
		{1_000, genderFemale, "тысяча", "тысячи", "тысяч"},
	}
	var parts []string
	for _, s := range scales {
		t := n / s.div
		n %= s.div
		if t > 0 {
			parts = append(parts, tripletToWords(t, s.gender), plural(t, s.one, s.few, s.many))
		}
	}
	if n > 0 {
		parts = append(parts, tripletToWords(n, gender))
	}
	return strings.Join(parts, " ")
}

func plural(n int64, one, few, many string) string {
	n = absInt64(n) % 100
	if n >= 11 && n <= 14 {
		return many
	}
	switch n % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return few
	default:
		return many
	}
}

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = []rune(strings.ToUpper(string(r[0])))[0]
	return string(r)
}

func SumField(items any, field string) float64 {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return 0
	}
	var sum float64
	for i := 0; i < v.Len(); i++ {
		it := v.Index(i)
		for it.Kind() == reflect.Ptr || it.Kind() == reflect.Interface {
			it = it.Elem()
		}
		if it.Kind() != reflect.Struct {
			continue
		}
		f := it.FieldByName(field)
		if !f.IsValid() {
			continue
		}
		switch f.Kind() {
		case reflect.Float32, reflect.Float64:
			sum += f.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sum += float64(f.Int())
		}
	}
	return sum
}
