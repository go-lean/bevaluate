package storage

import (
	"io"
	"os"
)

type (
	FileOpener struct{}
)

func (o FileOpener) OpenRead(path string) (io.ReadCloser, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModePerm)
}

func (o FileOpener) OpenCreate(path string) (io.WriteCloser, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}
