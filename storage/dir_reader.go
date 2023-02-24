package storage

import (
	"github.com/go-lean/bevaluate/models"
	"os"
)

type (
	DirReader struct{}

	DirEntry struct {
		name  string
		isDir bool
	}
)

func (r DirReader) Read(path string) ([]models.DirEntry, error) {
	entries, errRead := os.ReadDir(path)
	if errRead != nil {
		return nil, errRead
	}

	result := make([]models.DirEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, DirEntry{
			name:  entry.Name(),
			isDir: entry.IsDir(),
		})
	}

	return result, nil
}

func (e DirEntry) Name() string {
	return e.name
}

func (e DirEntry) IsDir() bool {
	return e.isDir
}
