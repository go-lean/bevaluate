package storage_test

import (
	"github.com/go-lean/bevaluate/storage"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOpener_OpenCreate(t *testing.T) {
	opener := storage.FileOpener{}
	path := filepath.Join(os.TempDir(), "baba.test")

	file, errCreate := opener.OpenCreate(path)
	require.NoError(t, errCreate)

	defer func() {
		_ = file.Close()
		_ = os.Remove(path)
	}()

	_, errWrite := io.WriteString(file, "baba is you")
	require.NoError(t, errWrite)

	data, errRead := os.ReadFile(path)

	require.NoError(t, errRead)
	require.Equal(t, "baba is you", string(data))
}

func TestFileOpener_OpenRead(t *testing.T) {
	path := filepath.Join(os.TempDir(), "baba.test")

	file, errCreate := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	require.NoError(t, errCreate)

	defer func() {
		_ = file.Close()
		_ = os.Remove(path)
	}()

	_, errWrite := io.WriteString(file, "baba is you")
	require.NoError(t, errWrite)

	opener := storage.FileOpener{}
	readFile, errRead := opener.OpenRead(path)
	require.NoError(t, errRead)

	defer func() {
		_ = readFile.Close()
	}()

	data, errRead := io.ReadAll(readFile)
	require.NoError(t, errRead)

	require.Equal(t, "baba is you", string(data))
}
