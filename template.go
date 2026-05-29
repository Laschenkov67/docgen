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
	}

	e.mu.Lock()
	e.templates[name] = &Template{Name: name, Raw: data, Format: f}
	e.mu.Unlock()
	return nil
}

func (e *Engine) LoadTemplateFile(name, path string) error {
	return e.LoadTemplate(name, FileSource{Path: path})
}

func (e *Engine) LoadTemplateFS(name string, fsys fs.FS, path string) error {
	return e.LoadTemplate(name, FSSource{FS: fsys, Path: path})
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
