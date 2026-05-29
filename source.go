package docgen

import (
	"io"
	"io/fs"
	"os"
)

type Source interface {
	Open() (io.ReadCloser, error)
	Name() string
}

type FileSource struct{ Path string }

func (f FileSource) Open() (io.ReadCloser, error) { return os.Open(f.Path) }
func (f FileSource) Name() string                 { return f.Path }

type FSSource struct {
	FS   fs.FS
	Path string
}

func (f FSSource) Open() (io.ReadCloser, error) { return f.FS.Open(f.Path) }
func (f FSSource) Name() string                 { return f.Path }

type BytesSource struct {
	N    string
	Data []byte
}

func (b BytesSource) Open() (io.ReadCloser, error) { return io.NopCloser(byteReader(b.Data)), nil }
func (b BytesSource) Name() string                 { return b.N }

type byteReader []byte

func (b byteReader) Read(p []byte) (int, error) {
	n := copy(p, b)
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}
