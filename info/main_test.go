package info_test

import (
	"errors"
	"github.com/go-lean/bevaluate/models"
	"io"
	"strings"
)

var (
	errKaboom      = errors.New("kaboom")
	testModuleName = "github.com/baba/is/you"
)

type (
	DataMock[T any] struct {
		items map[string]T
	}

	MockedDirReader struct {
		DataMock[[]models.DirEntry]
		canRead bool
	}

	MockedFileOpener struct {
		DataMock[io.ReadCloser]
		canOpen bool
	}

	DirEntry struct {
		name  string
		isDir bool
	}

	FakeFile struct {
		io.Reader
		canClose bool
	}
)

func NewDirReader() *MockedDirReader {
	return &MockedDirReader{
		DataMock: DataMock[[]models.DirEntry]{
			items: make(map[string][]models.DirEntry),
		},
		canRead: true,
	}
}

func NewFileOpener() *MockedFileOpener {
	return &MockedFileOpener{
		DataMock: DataMock[io.ReadCloser]{
			items: make(map[string]io.ReadCloser),
		},
		canOpen: true,
	}
}

func (r *DataMock[T]) MockAt(path string, data T) {
	if _, ok := r.items[path]; ok {
		panic("already mocked resource override")
	}

	r.items[path] = data
}

func (r *DataMock[T]) Read(path string) T {
	item, ok := r.items[path]
	if ok == false {
		panic("not mocked resource requested: " + path)
	}

	return item
}

func (r *MockedDirReader) Read(path string) ([]models.DirEntry, error) {
	if r.canRead == false {
		return nil, errKaboom
	}

	return r.DataMock.Read(path), nil
}

func (r *MockedFileOpener) OpenRead(path string) (io.ReadCloser, error) {
	if r.canOpen == false {
		return nil, errKaboom
	}

	return r.DataMock.Read(path), nil
}

func (r *MockedFileOpener) CanOpen(v bool) *MockedFileOpener {
	r.canOpen = v
	return r
}

func (r *MockedDirReader) CanRead(v bool) *MockedDirReader {
	r.canRead = v
	return r
}

func (e DirEntry) Name() string {
	return e.name
}

func (e DirEntry) IsDir() bool {
	return e.isDir
}

func NewFakeFile(content string) *FakeFile {
	return &FakeFile{
		canClose: true,
		Reader:   strings.NewReader(content),
	}
}

func (f *FakeFile) Close() error {
	if f.canClose == false {
		return errKaboom
	}

	return nil
}

func (f *FakeFile) CanClose(v bool) *FakeFile {
	f.canClose = v
	return f
}
