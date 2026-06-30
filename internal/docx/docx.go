package docx

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"strings"
	"text/template"

	"github.com/laschenkov67/docgen/internal/tmpl"
)

type Renderer struct {
	funcs  template.FuncMap
	strict bool
}

func New(funcs template.FuncMap, strict bool) *Renderer {
	return &Renderer{funcs: funcs, strict: strict}
}

var renderableParts = []string{
	"word/document.xml",
	"word/header1.xml",
	"word/header2.xml",
	"word/header3.xml",
	"word/footer1.xml",
	"word/footer2.xml",
	"word/footer3.xml",
}

func (r *Renderer) Render(ctx context.Context, tpl interface{ GetRaw() []byte }, data any, w io.Writer) error {
	return r.renderBytes(ctx, tpl.GetRaw(), data, w)
}

// Validate проверяет синтаксис всех рендерящихся частей .docx (без подстановки
// данных), чтобы битый шаблон обнаруживался при загрузке, а не при рендере.
func (r *Renderer) Validate(raw []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return fmt.Errorf("docx: open zip: %w", err)
	}
	for _, f := range zr.File {
		if !isRenderable(f.Name) {
			continue
		}
		if err := r.validatePart(f); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) validatePart(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("docx: open part %q: %w", f.Name, err)
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("docx: read %q: %w", f.Name, err)
	}
	if err := tmpl.Validate(f.Name, string(normalizeForTemplate(body)), r.funcs, r.strict); err != nil {
		return fmt.Errorf("docx: %q: %w", f.Name, err)
	}
	return nil
}

func (r *Renderer) renderBytes(ctx context.Context, raw []byte, data any, w io.Writer) error {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return fmt.Errorf("docx: open zip: %w", err)
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	escaped := escapeTemplateData(data)

	for _, f := range zr.File {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := r.copyOrRender(f, zw, escaped); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) copyOrRender(f *zip.File, zw *zip.Writer, data any) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("docx: open part %q: %w", f.Name, err)
	}
	defer rc.Close()

	out, err := zw.CreateHeader(&zip.FileHeader{
		Name:   f.Name,
		Method: f.Method,
	})
	if err != nil {
		return fmt.Errorf("docx: create header %q: %w", f.Name, err)
	}

	if !isRenderable(f.Name) {
		_, err = io.Copy(out, rc)
		return err
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("docx: read %q: %w", f.Name, err)
	}

	normalized := normalizeForTemplate(body)

	rendered, err := tmpl.Execute(f.Name, string(normalized), r.funcs, data, r.strict)
	if err != nil {
		return fmt.Errorf("docx: render %q: %w", f.Name, err)
	}

	_, err = out.Write(rendered)
	return err
}

func normalizeForTemplate(body []byte) []byte {
	normalized := NormalizeRuns(body)
	return decodeEntitiesInTags(normalized)
}

func isRenderable(name string) bool {
	for _, p := range renderableParts {
		if p == name {
			return true
		}
	}
	return false
}

func decodeEntitiesInTags(src []byte) []byte {
	s := string(src)
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		j := strings.Index(s[i:], "{{")
		if j < 0 {
			b.WriteString(s[i:])
			break
		}
		b.WriteString(s[i : i+j])
		k := strings.Index(s[i+j:], "}}")
		if k < 0 {
			b.WriteString(s[i+j:])
			break
		}
		tag := s[i+j : i+j+k+2]
		b.WriteString(html.UnescapeString(tag))
		i += j + k + 2
	}
	return []byte(b.String())
}
