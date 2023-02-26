package storage

import (
	"errors"
	"io"
	"os"
)

var ErrNotExisting = errors.New("path does not exist")

type (
	FileReadOpener interface {
		OpenRead(path string) (io.ReadCloser, error)
	}

	FileCreateOpener interface {
		OpenCreate(path string) (io.WriteCloser, error)
	}

	Store struct {
		FileOpener
		DirReader
	}
)

func (s Store) TryAccessing(path string) error {
	_, errStat := os.Stat(path)
	if os.IsNotExist(errStat) {
		return ErrNotExisting
	}

	return errStat
}
