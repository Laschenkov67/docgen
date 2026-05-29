package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/laschenkov67/docgen"
)

func main() {
	eng := docgen.New()
	src := []byte(`#TITLE Акт оказанных услуг №{{.N}} от {{dateLong .Date}}
#H1 Исполнитель
{{.Performer}}
#H1 Заказчик
{{.Customer}}

Стороны подтверждают, что услуги оказаны в полном объёме.
Стоимость услуг: {{money .Amount}} руб. ({{rubles .Amount}}).
`)
	_ = eng.LoadTemplateBytes("act", src, docgen.FormatPDF)

	f, _ := os.Create("act.pdf")
	defer f.Close()
	_ = eng.Render(context.Background(), "act", map[string]any{
		"N": "7", "Date": time.Now(),
		"Performer": "ООО «АгроСофт»",
		"Customer":  "КФХ Иванов И.И.",
		"Amount":    49999.99,
	}, docgen.FormatPDF, f)
	log.Println("ok: act.pdf")
}
