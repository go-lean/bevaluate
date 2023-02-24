package storage_test

import (
	"errors"
	"io"
)

type (
	FakeReadCloser struct {
		io.Reader
		canRead  bool
		canClose bool
	}

	FakeWriteCloser struct {
		io.Writer
		canWrite bool
		canClose bool
	}
)

var (
	errKaboom = errors.New("kaboom")
)

func NewFakeReadCloser(canRead, canClose bool, reader io.Reader) FakeReadCloser {
	return FakeReadCloser{
		Reader:   reader,
		canRead:  canRead,
		canClose: canClose,
	}
}

func (r FakeReadCloser) Read(b []byte) (int, error) {
	if r.canRead == false {
		return 0, errKaboom
	}

	return r.Reader.Read(b)
}

func (r FakeReadCloser) Close() error {
	if r.canClose == false {
		return errKaboom
	}

	return nil
}

func NewFakeWriteCloser(canWrite, canClose bool, writer io.Writer) FakeWriteCloser {
	return FakeWriteCloser{
		Writer:   writer,
		canWrite: canWrite,
		canClose: canClose,
	}
}

func (r FakeWriteCloser) Write(b []byte) (int, error) {
	if r.canWrite == false {
		return 0, errKaboom
	}

	return r.Writer.Write(b)
}

func (r FakeWriteCloser) Close() error {
	if r.canClose == false {
		return errKaboom
	}

	return nil
}
