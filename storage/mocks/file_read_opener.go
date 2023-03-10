// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// FileReadOpener is an autogenerated mock type for the FileReadOpener type
type FileReadOpener struct {
	mock.Mock
}

// OpenReadOnly provides a mock function with given fields: path
func (_m *FileReadOpener) OpenRead(path string) (io.ReadCloser, error) {
	ret := _m.Called(path)

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (io.ReadCloser, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFileReadOpener interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileReadOpener creates a new instance of FileReadOpener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileReadOpener(t mockConstructorTestingTNewFileReadOpener) *FileReadOpener {
	mock := &FileReadOpener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
