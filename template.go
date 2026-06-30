package docgen

import (
	"fmt"
	"io"
	"io/fs"
)

type Template struct {
	Name   string
	Raw    []byte
	Format Format
}

// GetRaw реализует минимальный контракт, который ожидают Renderer-ы
// (см. renderer.go) — отдаёт исходные байты шаблона.
func (t *Template) GetRaw() []byte { return t.Raw }

func (e *Engine) LoadTemplate(name string, src Source, format ...Format) error {
	rc, err := src.Open()
	if err != nil {
		return fmt.Errorf("docgen: open template %q: %w", name, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("docgen: read template %q: %w", name, err)
	}

	var f Format
	if len(format) > 0 {
		f = format[0]
	} else {
		f = detectFormat(src.Name())
		if f == "" {
			return fmt.Errorf("%w: %q: не удалось определить формат по расширению %q, укажите формат явно", ErrInvalidTemplate, name, src.Name())
		}
	}

	e.mu.RLock()
	r, rok := e.renderers[f]
	e.mu.RUnlock()
	if rok {
		if v, ok := r.(interface{ Validate(raw []byte) error }); ok {
			if err := v.Validate(data); err != nil {
				return fmt.Errorf("%w: %q: %w", ErrInvalidTemplate, name, err)
			}
		}
	}

	e.mu.Lock()
	e.templates[name] = &Template{Name: name, Raw: data, Format: f}
	e.mu.Unlock()
	return nil
}

func (e *Engine) LoadTemplateFile(name, path string, format ...Format) error {
	return e.LoadTemplate(name, FileSource{Path: path}, format...)
}

func (e *Engine) LoadTemplateFS(name string, fsys fs.FS, path string, format ...Format) error {
	return e.LoadTemplate(name, FSSource{FS: fsys, Path: path}, format...)
}

func (e *Engine) LoadTemplateBytes(name string, data []byte, format Format) error {
	return e.LoadTemplate(name, BytesSource{N: name, Data: data}, format)
}

func detectFormat(name string) Format {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			switch name[i+1:] {
			case "docx":
				return FormatDOCX
			case "pdf":
				return FormatPDF
			}
			break
		}
	}
	return ""
}
