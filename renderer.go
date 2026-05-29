package docgen

import (
	"context"
	"io"
)

type Renderer interface {
	Render(ctx context.Context, tpl *Template, data any, w io.Writer) error
}

type RendererFunc func(ctx context.Context, tpl *Template, data any, w io.Writer) error

func (f RendererFunc) Render(ctx context.Context, tpl *Template, data any, w io.Writer) error {
	return f(ctx, tpl, data, w)
}
