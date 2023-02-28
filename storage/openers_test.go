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

func TestFileOpener_OpenCreate_DirectoryNotExisting_DirectoryIsCreated(t *testing.T) {
	opener := storage.FileOpener{}
	path := filepath.Join(os.TempDir(), "some", "inner", "folder", "baba.test")

	file, errCreate := opener.OpenCreate(path)
	require.NoError(t, errCreate)

	defer func() {
		_ = file.Close()
		_ = os.RemoveAll(filepath.Join(os.TempDir(), "some"))
	}()

	_, errWrite := io.WriteString(file, "baba is you")
	require.NoError(t, errWrite)

	data, errRead := os.ReadFile(path)

	require.NoError(t, errRead)
	require.Equal(t, "baba is you", string(data))
}

func TestFileOpener_OpenCreate_DirectoryExists_NoError(t *testing.T) {
	opener := storage.FileOpener{}
	path := filepath.Join(os.TempDir(), "some", "inner", "folder", "baba.test")
	errDir := os.MkdirAll(filepath.Join(os.TempDir(), "some", "inner", "folder"), os.ModePerm)
	require.NoError(t, errDir)
	defer func() {
		_ = os.RemoveAll(filepath.Join(os.TempDir(), "some"))
	}()

	file, errCreate := opener.OpenCreate(path)
	require.NoError(t, errCreate)

	defer func() {
		_ = file.Close()
	}()

	_, errWrite := io.WriteString(file, "baba is you")
	require.NoError(t, errWrite)

	data, errRead := os.ReadFile(path)

	require.NoError(t, errRead)
	require.Equal(t, "baba is you", string(data))
}

func TestFileOpener_OpenCreate_DirectoryNotExisting_FileWithSameName_Error(t *testing.T) {
	opener := storage.FileOpener{}
	path := filepath.Join(os.TempDir(), "some", "inner", "folder", "baba.test")
	errDir := os.MkdirAll(filepath.Join(os.TempDir(), "some", "inner"), os.ModePerm)
	require.NoError(t, errDir)
	defer func() {
		_ = os.RemoveAll(filepath.Join(os.TempDir(), "some"))
	}()

	preFile, preErrCreate := os.OpenFile(filepath.Join(os.TempDir(), "some", "inner", "folder"), os.O_CREATE, os.ModePerm)
	require.NoError(t, preErrCreate)
	require.NoError(t, preFile.Close())

	file, errCreate := opener.OpenCreate(path)
	require.Error(t, errCreate)
	require.Nil(t, file)
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
