package docgen

import (
	"context"
	"io"
)

// Renderer-у достаточно прочитать сырые байты шаблона: сигнатура использует
// анонимный интерфейс (а не *Template), чтобы internal/docx и internal/pdf
// могли реализовать его, не импортируя пакет docgen (иначе возник бы цикл
// docgen -> internal/docx -> docgen).
type Renderer interface {
	Render(ctx context.Context, tpl interface{ GetRaw() []byte }, data any, w io.Writer) error
}

type RendererFunc func(ctx context.Context, tpl interface{ GetRaw() []byte }, data any, w io.Writer) error

func (f RendererFunc) Render(ctx context.Context, tpl interface{ GetRaw() []byte }, data any, w io.Writer) error {
	return f(ctx, tpl, data, w)
}
