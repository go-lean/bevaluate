package storage_test

import (
	"github.com/go-lean/bevaluate/storage"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestDirReader_Read_NonExistingPath_Error(t *testing.T) {
	reader := storage.DirReader{}
	entries, errRead := reader.Read("baba")

	require.Error(t, errRead)
	require.Empty(t, entries)
}

func TestDirReader_Read_OK(t *testing.T) {
	path := filepath.Join(os.TempDir(), "baba")
	errMakeDir := os.Mkdir(path, os.ModePerm)
	require.NoError(t, errMakeDir)
	defer func() {
		_ = os.RemoveAll(path)
	}()

	errMakeDir = os.Mkdir(filepath.Join(path, "subDir"), os.ModePerm)
	require.NoError(t, errMakeDir)

	f, errCreate := os.Create(filepath.Join(path, "baba.go"))
	require.NoError(t, errCreate)
	require.NoError(t, f.Close())

	reader := storage.DirReader{}
	entries, errRead := reader.Read(path)

	require.NoError(t, errRead)
	require.Len(t, entries, 2)

	require.Equal(t, "baba.go", entries[0].Name())
	require.False(t, entries[0].IsDir())

	require.Equal(t, "subDir", entries[1].Name())
	require.True(t, entries[1].IsDir())
}
