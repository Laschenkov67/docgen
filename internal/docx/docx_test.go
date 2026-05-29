package docx_test

import (
	"archive/zip"
	"bytes"
	"context"
	"testing"

	"github.com/laschenkov67/docgen/internal/docx"
	"github.com/laschenkov67/docgen/internal/tmpl"
	"github.com/stretchr/testify/require"
)

// buildMinimalDocx собирает минимальный валидный .docx с указанным document.xml.
func buildMinimalDocx(t *testing.T, documentXML string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="xml" ContentType="application/xml"/>
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		"word/document.xml": documentXML,
	}
	for name, body := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(body))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

type rawTpl struct{ raw []byte }

func (r rawTpl) GetRaw() []byte { return r.raw }

func TestRender_Simple(t *testing.T) {
	doc := buildMinimalDocx(t, `<?xml version="1.0"?><w:document xmlns:w="x"><w:body>
<w:p><w:r><w:t>Здравствуйте, {{.Name}}!</w:t></w:r></w:p>
</w:body></w:document>`)

	r := docx.New(tmpl.DefaultFuncs(), true)
	var out bytes.Buffer
	err := r.Render(context.Background(), rawTpl{doc}, map[string]any{"Name": "Иван"}, &out)
	require.NoError(t, err)

	body, err := docx.ExtractPart(out.Bytes(), "word/document.xml")
	require.NoError(t, err)
	require.Contains(t, string(body), "Здравствуйте, Иван!")
}

func TestRender_SplitRuns(t *testing.T) {
	// Плейсхолдер разорван между двумя runs — типичный кейс Word'а.
	doc := buildMinimalDocx(t, `<?xml version="1.0"?><w:document xmlns:w="x"><w:body>
<w:p>
  <w:r><w:t>Привет, {{ .Na</w:t></w:r>
  <w:r><w:t>me }}!</w:t></w:r>
</w:p>
</w:body></w:document>`)

	r := docx.New(tmpl.DefaultFuncs(), true)
	var out bytes.Buffer
	err := r.Render(context.Background(), rawTpl{doc}, map[string]any{"Name": "Мир"}, &out)
	require.NoError(t, err)

	body, err := docx.ExtractPart(out.Bytes(), "word/document.xml")
	require.NoError(t, err)
	require.Contains(t, string(body), "Привет, Мир!")
}

func TestRender_Loop(t *testing.T) {
	doc := buildMinimalDocx(t, `<?xml version="1.0"?><w:document xmlns:w="x"><w:body>
<w:p><w:r><w:t>{{range .Items}}{{.Name}}={{.Qty}};{{end}}</w:t></w:r></w:p>
</w:body></w:document>`)

	type Item struct {
		Name string
		Qty  int
	}
	data := map[string]any{
		"Items": []Item{{"A", 1}, {"B", 2}},
	}

	r := docx.New(tmpl.DefaultFuncs(), true)
	var out bytes.Buffer
	require.NoError(t, r.Render(context.Background(), rawTpl{doc}, data, &out))

	body, _ := docx.ExtractPart(out.Bytes(), "word/document.xml")
	require.Contains(t, string(body), "A=1;B=2;")
}
