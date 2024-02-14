package config

import (
	"bytes"
	"io"

	"github.com/stretchr/testify/mock"
)

type __rwMockArgs struct {
	// Write and Read errors
	writeErr error
	readErr  error

	// __RWC errors
	writerErr error // io.WriteCloser error on Write
	readerErr error // io.ReadCloser error on Read
	closerErr error // io.ReadWriteCloser error on Close

	// Exists
	exists bool

	// data init for __RWC.buf
	data []byte
}

func __newRWMock(a __rwMockArgs) *__rwMock {
	r := &__rwMock{
		rwc: &__RWC{},
	}

	r.rwc.writeError = a.writerErr
	r.rwc.closeError = a.closerErr
	r.rwc.readError = a.readerErr

	if len(a.data) != 0 {
		r.rwc.buf.Write(a.data)
	}
	r.On("Write", mock.Anything, mock.Anything).Return(r.rwc, a.writeErr)
	r.On("Read", mock.Anything).Return(r.rwc, a.readErr)
	r.On("Exists").Return(a.exists)

	return r
}

type __rwMock struct {
	mock.Mock
	rwc *__RWC
}

func (rw *__rwMock) Write(patch string, flag int) (io.WriteCloser, error) {
	ret := rw.Called(patch, flag)

	return ret.Get(0).(io.WriteCloser), ret.Error(1)
}

func (rw *__rwMock) Read(patch string) (io.ReadCloser, error) {
	ret := rw.Called(patch)

	return ret.Get(0).(io.ReadCloser), ret.Error(1)
}

func (rw *__rwMock) Exists(_ string) bool {
	ret := rw.Called()

	return ret.Bool(0)
}

// __RWC - simple structure io.ReadWriteCloser.
type __RWC struct {
	buf        bytes.Buffer
	writeError error
	readError  error
	closeError error
}

func (rw *__RWC) Write(p []byte) (n int, err error) {
	if rw.writeError != nil {
		return 0, rw.writeError
	}

	return rw.buf.Write(p)
}

func (rw *__RWC) Read(p []byte) (n int, err error) {
	if rw.readError != nil {
		return 0, rw.readError
	}

	return rw.buf.Read(p)
}

func (rw *__RWC) Close() error {
	return rw.closeError
}
