package format_test

import (
	"testing"

	"github.com/laschenkov67/docgen/format"
	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	cases := map[float64]string{
		0:          "0,00",
		1:          "1,00",
		1234.5:     "1 234,50",
		1234567.89: "1 234 567,89",
		-99.999:    "-100,00",
	}
	for in, want := range cases {
		assert.Equal(t, want, format.Money(in), "v=%v", in)
	}
}

func TestRubles(t *testing.T) {
	cases := map[float64]string{
		0:       "Ноль рублей 00 копеек",
		1:       "Один рубль 00 копеек",
		2.5:     "Два рубля 50 копеек",
		21:      "Двадцать один рубль 00 копеек",
		1001:    "Одна тысяча один рубль 00 копеек",
		1234.56: "Одна тысяча двести тридцать четыре рубля 56 копеек",
	}
	for in, want := range cases {
		assert.Equal(t, want, format.Rubles(in), "v=%v", in)
	}
}

func TestSumField(t *testing.T) {
	type item struct{ Total float64 }
	items := []item{{10.5}, {20.25}, {69.25}}
	assert.InDelta(t, 100.0, format.SumField(items, "Total"), 0.0001)
}
