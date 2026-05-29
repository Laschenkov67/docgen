package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

func Execute(name, src string, funcs template.FuncMap, data any, strict bool) ([]byte, error) {
	t := template.New(name).Funcs(funcs)
	if strict {
		t = t.Option("missingkey=error")
	} else {
		t = t.Option("missingkey=zero")
	}
	t, err := t.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}
	return buf.Bytes(), nil
}
