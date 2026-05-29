package docgen

import (
	"text/template"

	"github.com/laschenkov67/docgen/internal/docx"
	"github.com/laschenkov67/docgen/internal/pdf"
	"github.com/laschenkov67/docgen/internal/tmpl"
)

type Option func(*options)

type options struct {
	funcs     template.FuncMap
	renderers map[Format]Renderer
	strict    bool // missingkey=error для text/template
}

func defaultOptions() options {
	o := options{
		funcs:     tmpl.DefaultFuncs(),
		renderers: map[Format]Renderer{},
		strict:    true,
	}
	o.renderers[FormatDOCX] = docx.New(o.funcs, o.strict)
	o.renderers[FormatPDF] = pdf.New(o.funcs, o.strict)
	return o
}

func WithFuncs(fns template.FuncMap) Option {
	return func(o *options) {
		for k, v := range fns {
			o.funcs[k] = v
		}
		o.renderers[FormatDOCX] = docx.New(o.funcs, o.strict)
		o.renderers[FormatPDF] = pdf.New(o.funcs, o.strict)
	}
}

func WithRenderer(format Format, r Renderer) Option {
	return func(o *options) {
		o.renderers[format] = r
	}
}

func WithStrict(strict bool) Option {
	return func(o *options) {
		o.strict = strict
		o.renderers[FormatDOCX] = docx.New(o.funcs, o.strict)
		o.renderers[FormatPDF] = pdf.New(o.funcs, o.strict)
	}
}
