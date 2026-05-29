package docgen

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type Format string

const (
	FormatDOCX Format = "docx"
	FormatPDF  Format = "pdf"
)

type Engine struct {
	mu        sync.RWMutex
	templates map[string]*Template
	renderers map[Format]Renderer
	opts      options
}

func New(opts ...Option) *Engine {
	o := defaultOptions()
	for _, fn := range opts {
		fn(&o)
	}
	e := &Engine{
		templates: make(map[string]*Template),
		renderers: make(map[Format]Renderer, len(o.renderers)),
		opts:      o,
	}
	for f, r := range o.renderers {
		e.renderers[f] = r
	}
	return e
}

func (e *Engine) Render(ctx context.Context, name string, data any, format Format, w io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	e.mu.RLock()
	tpl, ok := e.templates[name]
	r, rok := e.renderers[format]
	e.mu.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %q", ErrTemplateNotFound, name)
	}
	if !rok {
		return fmt.Errorf("%w: %q", ErrRendererNotFound, format)
	}

	return r.Render(ctx, tpl, data, w)
}

func (e *Engine) HasTemplate(name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.templates[name]
	return ok
}

func (e *Engine) Templates() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]string, 0, len(e.templates))
	for n := range e.templates {
		out = append(out, n)
	}
	return out
}
