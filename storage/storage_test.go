package storage_test

import (
	"github.com/go-lean/bevaluate/storage"
	"github.com/go-lean/bevaluate/storage/mocks"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// region Read Module Name

func TestReadModuleName_OpenError(t *testing.T) {
	opener := mocks.NewFileReadOpener(t)
	opener.On("OpenRead", "baba/go.mod").
		Return(nil, errKaboom)

	name, err := storage.ReadModuleName("baba/go.mod", opener)

	require.Error(t, err)
	require.Contains(t, err.Error(), "kaboom")
	require.NotEqual(t, err.Error(), "kaboom")
	require.Empty(t, name)
}

func TestReadModuleName_ReadError(t *testing.T) {
	reader := NewFakeReadCloser(false, true, nil)
	opener := mocks.NewFileReadOpener(t)
	opener.On("OpenRead", "baba/go.mod").
		Return(reader, nil)

	name, err := storage.ReadModuleName("baba/go.mod", opener)

	require.Error(t, err)
	require.Contains(t, err.Error(), "kaboom")
	require.NotEqual(t, err.Error(), "kaboom")
	require.Empty(t, name)
}

func TestReadModuleName_CloseError(t *testing.T) {
	r := strings.NewReader("module github.com/baba/is/you")
	reader := NewFakeReadCloser(true, false, r)
	opener := mocks.NewFileReadOpener(t)
	opener.On("OpenRead", "baba/go.mod").
		Return(reader, nil)

	name, err := storage.ReadModuleName("baba/go.mod", opener)

	require.NoError(t, err)
	require.Equal(t, "github.com/baba/is/you", name)
}

func TestReadModuleName_NewLine_ShouldOnlyGetName(t *testing.T) {
	r := strings.NewReader("module github.com/baba/is/you\n\nsome other stuff")
	reader := NewFakeReadCloser(true, true, r)
	opener := mocks.NewFileReadOpener(t)
	opener.On("OpenRead", "baba/go.mod").
		Return(reader, nil)

	name, err := storage.ReadModuleName("baba/go.mod", opener)

	require.NoError(t, err)
	require.Equal(t, "github.com/baba/is/you", name)
}

// endregion Read Module Name

// region Create File With Text

func TestCreateFileWithText_OpenError(t *testing.T) {
	opener := mocks.NewFileCreateOpener(t)
	opener.On("OpenCreate", "baba.go").
		Return(nil, errKaboom)

	err := storage.CreateFileWithText("baba.go", "baba is you", opener)

	require.Error(t, err)
	require.Contains(t, err.Error(), "kaboom")
	require.NotEqual(t, err.Error(), "kaboom")
}

func TestCreateFileWithText_WriteError(t *testing.T) {
	writer := NewFakeWriteCloser(false, true, nil)
	opener := mocks.NewFileCreateOpener(t)
	opener.On("OpenCreate", "baba.go").
		Return(writer, nil)

	err := storage.CreateFileWithText("baba.go", "baba is you", opener)

	require.Error(t, err)
	require.Contains(t, err.Error(), "kaboom")
	require.NotEqual(t, err.Error(), "kaboom")
}

func TestCreateFileWithText_CloseError_Continue(t *testing.T) {
	sb := strings.Builder{}
	writer := NewFakeWriteCloser(true, false, &sb)
	opener := mocks.NewFileCreateOpener(t)
	opener.On("OpenCreate", "baba.go").
		Return(writer, nil)

	err := storage.CreateFileWithText("baba.go", "baba is you", opener)

	require.NoError(t, err)
	require.Equal(t, "baba is you", sb.String())
}

func TestCreateFileWithText_OK(t *testing.T) {
	sb := strings.Builder{}
	writer := NewFakeWriteCloser(true, true, &sb)
	opener := mocks.NewFileCreateOpener(t)
	opener.On("OpenCreate", "baba.go").
		Return(writer, nil)

	err := storage.CreateFileWithText("baba.go", "baba is you", opener)

	require.NoError(t, err)
	require.Equal(t, "baba is you", sb.String())
}

// endregion Create File With Text
