package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

func Execute(name, src string, funcs template.FuncMap, data any, strict bool) ([]byte, error) {
	t, err := parse(name, src, funcs, strict)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}
	return buf.Bytes(), nil
}

// Validate проверяет только синтаксис шаблона (без выполнения), чтобы ошибку
// в разметке можно было поймать на этапе загрузки шаблона, а не во время
// рендера у конечного пользователя.
func Validate(name, src string, funcs template.FuncMap, strict bool) error {
	_, err := parse(name, src, funcs, strict)
	return err
}

func parse(name, src string, funcs template.FuncMap, strict bool) (*template.Template, error) {
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
	return t, nil
}
