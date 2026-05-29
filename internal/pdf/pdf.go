package pdf

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/go-pdf/fpdf"
	"github.com/laschenkov67/docgen/internal/tmpl"
)

// Renderer для PDF.
type Renderer struct {
	funcs  template.FuncMap
	strict bool
}

// New создаёт PDF renderer.
func New(funcs template.FuncMap, strict bool) *Renderer {
	return &Renderer{funcs: funcs, strict: strict}
}

// Render — реализация docgen.Renderer.
func (r *Renderer) Render(ctx context.Context, tpl interface{ GetRaw() []byte }, data any, w io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	rendered, err := tmpl.Execute("pdf", string(tpl.GetRaw()), r.funcs, data, r.strict)
	if err != nil {
		return fmt.Errorf("pdf: render template: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("DejaVu", "", dejaVuRegular)
	pdf.AddUTF8FontFromBytes("DejaVu", "B", dejaVuBold)
	pdf.SetFont("DejaVu", "", 11)
	pdf.AddPage()

	for _, line := range strings.Split(string(rendered), "\n") {
		line = strings.TrimRight(line, "\r ")
		if err := emitLine(pdf, line); err != nil {
			return err
		}
	}

	if err := pdf.Output(w); err != nil {
		return fmt.Errorf("pdf: output: %w", err)
	}
	return nil
}

func emitLine(pdf *fpdf.Fpdf, line string) error {
	switch {
	case strings.HasPrefix(line, "#TITLE "):
		pdf.SetFont("DejaVu", "B", 16)
		pdf.MultiCell(0, 10, strings.TrimPrefix(line, "#TITLE "), "", "C", false)
		pdf.Ln(2)
		pdf.SetFont("DejaVu", "", 11)
	case strings.HasPrefix(line, "#H1 "):
		pdf.SetFont("DejaVu", "B", 13)
		pdf.MultiCell(0, 8, strings.TrimPrefix(line, "#H1 "), "", "L", false)
		pdf.SetFont("DejaVu", "", 11)
	case strings.HasPrefix(line, "#TABLE "):
		headers := strings.Split(strings.TrimPrefix(line, "#TABLE "), "|")
		emitTableHeader(pdf, headers)
	case strings.Contains(line, "|") && !strings.HasPrefix(line, "#"):
		cells := strings.Split(line, "|")
		emitTableRow(pdf, cells)
	case strings.HasPrefix(line, "#TOTAL "):
		pdf.Ln(2)
		pdf.SetFont("DejaVu", "B", 12)
		pdf.MultiCell(0, 8, strings.TrimPrefix(line, "#TOTAL "), "", "R", false)
		pdf.SetFont("DejaVu", "", 11)
	case line == "":
		pdf.Ln(4)
	default:
		pdf.MultiCell(0, 6, line, "", "L", false)
	}
	return nil
}

func emitTableHeader(pdf *fpdf.Fpdf, headers []string) {
	pdf.SetFont("DejaVu", "B", 11)
	colW := tableColWidths(pdf, len(headers))
	for i, h := range headers {
		pdf.CellFormat(colW[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("DejaVu", "", 11)
}

func emitTableRow(pdf *fpdf.Fpdf, cells []string) {
	colW := tableColWidths(pdf, len(cells))
	for i, c := range cells {
		pdf.CellFormat(colW[i], 7, c, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)
}

func tableColWidths(pdf *fpdf.Fpdf, n int) []float64 {
	w, _ := pdf.GetPageSize()
	l, _, r, _ := pdf.GetMargins()
	avail := w - l - r
	col := avail / float64(n)
	out := make([]float64, n)
	for i := range out {
		out[i] = col
	}
	return out
}
