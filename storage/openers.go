package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type (
	FileOpener struct{}
)

func (o FileOpener) OpenRead(path string) (io.ReadCloser, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModePerm)
}

func (o FileOpener) OpenCreate(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	if errMkDir := os.MkdirAll(dir, os.ModePerm); errMkDir != nil {
		return nil, fmt.Errorf("could not create directory for new file: %w", errMkDir)
	}

	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}
