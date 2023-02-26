package storage_test

import (
	"github.com/go-lean/bevaluate/storage"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestStore_Exists_NonExistingPath_False(t *testing.T) {
	s := storage.Store{}
	err := s.TryAccessing("baba.go")

	require.Error(t, err)
	require.ErrorIs(t, err, storage.ErrNotExisting)
}

func TestStore_Exists_ExistingPath_True(t *testing.T) {
	path := filepath.Join(os.TempDir(), "baba.go")
	file, errCreate := os.Create(path)
	require.NoError(t, errCreate)

	_ = file.Close()
	defer func() {
		_ = os.Remove(path)
	}()

	s := storage.Store{}
	err := s.TryAccessing(path)

	require.NoError(t, err)
}
