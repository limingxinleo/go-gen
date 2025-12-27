package module

import (
	"os"

	"golang.org/x/mod/modfile"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "github.com/limingxinleo/go-zero-skeleton"
	}

	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "github.com/limingxinleo/go-zero-skeleton"
	}

	return file.Module.Mod.Path
}
