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

func (r *Renderer) renderBytes(ctx context.Context, raw []byte, data any, w io.Writer) error {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return fmt.Errorf("docx: open zip: %w", err)
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	for _, f := range zr.File {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := r.copyOrRender(f, zw, data); err != nil {
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

	normalized := NormalizeRuns(body)

	normalized = decodeEntitiesInTags(normalized)

	rendered, err := tmpl.Execute(f.Name, string(normalized), r.funcs, data, r.strict)
	if err != nil {
		return fmt.Errorf("docx: render %q: %w", f.Name, err)
	}

	_, err = out.Write(rendered)
	return err
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
