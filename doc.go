// Package docgen — универсальный генератор документов (DOCX, PDF) из шаблонов.
//
// Поддерживаются:
//
//   - DOCX-шаблоны с плейсхолдерами синтаксиса text/template ({{.Customer.Name}}).
//     Корректно обрабатываются разорванные текстовые runs Word.
//   - PDF-генерация через gofpdf (для простых табличных документов:
//     накладные, акты, счета).
//   - Pluggable Renderer-интерфейс: можно добавить свой бэкенд
//     (HTML, ODT, wkhtmltopdf, chromedp и т. п.) через docgen.RegisterRenderer.
//   - FuncMap: rubles, date, money, inn, kpp, upper, lower, default, и т. д.
//
// Пример:
//
//	eng := docgen.New()
//	_ = eng.LoadTemplateFile("invoice", "templates/invoice.docx")
//
//	var buf bytes.Buffer
//	err := eng.Render(ctx, "invoice", data, docgen.FormatDOCX, &buf)
//
// Подробнее см. examples/.
package docgen
