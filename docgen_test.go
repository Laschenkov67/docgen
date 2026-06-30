package docgen_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/laschenkov67/docgen"
	"github.com/stretchr/testify/require"
)

func TestEngine_PDF(t *testing.T) {
	eng := docgen.New()
	src := []byte(`#TITLE Привет, {{.Name}}!`)
	require.NoError(t, eng.LoadTemplateBytes("hello", src, docgen.FormatPDF))

	var out bytes.Buffer
	err := eng.Render(context.Background(), "hello", map[string]any{"Name": "мир"}, docgen.FormatPDF, &out)
	require.NoError(t, err)
	require.True(t, bytes.HasPrefix(out.Bytes(), []byte("%PDF-")))
}

func TestEngine_TemplateNotFound(t *testing.T) {
	eng := docgen.New()
	err := eng.Render(context.Background(), "nope", nil, docgen.FormatPDF, &bytes.Buffer{})
	require.ErrorIs(t, err, docgen.ErrTemplateNotFound)
}

func TestEngine_StrictMode(t *testing.T) {
	eng := docgen.New(docgen.WithStrict(true))
	require.NoError(t, eng.LoadTemplateBytes("t", []byte(`{{.Missing}}`), docgen.FormatPDF))
	err := eng.Render(context.Background(), "t", map[string]any{}, docgen.FormatPDF, &bytes.Buffer{})
	require.Error(t, err)
	require.ErrorIs(t, err, docgen.ErrRenderFailed)
}

func TestEngine_LoadTemplate_InvalidSyntax(t *testing.T) {
	eng := docgen.New()
	err := eng.LoadTemplateBytes("broken", []byte(`{{.Unclosed`), docgen.FormatPDF)
	require.ErrorIs(t, err, docgen.ErrInvalidTemplate)
}
