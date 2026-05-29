package pdf_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/laschenkov67/docgen/internal/pdf"
	"github.com/laschenkov67/docgen/internal/tmpl"
	"github.com/stretchr/testify/require"
)

type rawTpl struct{ raw []byte }

func (r rawTpl) GetRaw() []byte { return r.raw }

func TestPDF_Smoke(t *testing.T) {
	src := `#TITLE Счёт №{{.N}}
#H1 Поставщик
ООО «Ромашка», ИНН 7707083893
#TABLE Наименование|Кол-во|Цена|Сумма
{{range .Items}}{{.Name}}|{{.Qty}}|{{money .Price}}|{{money .Total}}
{{end}}#TOTAL Итого: {{money (sum .Items "Total")}}
`
	type Item struct {
		Name         string
		Qty          int
		Price, Total float64
	}
	data := map[string]any{
		"N": "42",
		"Items": []Item{
			{"Картофель, кг", 10, 35.5, 355},
			{"Морковь, кг", 5, 42, 210},
		},
	}

	r := pdf.New(tmpl.DefaultFuncs(), true)
	var out bytes.Buffer
	require.NoError(t, r.Render(context.Background(), rawTpl{[]byte(src)}, data, &out))
	require.True(t, bytes.HasPrefix(out.Bytes(), []byte("%PDF-")))
}
