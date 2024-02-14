package backup

import (
	"bytes"
	"io"

	"github.com/klimby/version/pkg/convert"
	"github.com/stretchr/testify/mock"
)

const __patch = "source"

type __rwMockArgs struct {
	// Write and Read errors
	writeErr error
	readErr  error

	// __RWC errors
	writerErr error // io.WriteCloser error on Write
	readerErr error // io.ReadCloser error on Read
	closerErr error // io.ReadWriteCloser error on Close

	// RemoveAll
	removeErr error

	// data init for __RWC.buf
	data []byte
}

func __newRWMock(source __rwMockArgs, target __rwMockArgs) *__rwMock {
	r := &__rwMock{
		rwcSource: &__RWC{},
		rwcTarget: &__RWC{},
	}

	r.rwcSource.writeError = source.writerErr
	r.rwcSource.closeError = source.closerErr
	r.rwcSource.readError = source.readerErr

	r.rwcTarget.writeError = target.writerErr
	r.rwcTarget.closeError = target.closerErr
	r.rwcTarget.readError = target.readerErr

	if len(source.data) != 0 {
		r.rwcSource.buf.Write(source.data)
	}

	if len(target.data) != 0 {
		r.rwcTarget.buf.Write(target.data)
	}

	r.On("Read", __patch).Return(r.rwcSource, source.readErr)

	backPath := __patch + _suffix
	r.On("Write", backPath, mock.Anything).Return(r.rwcTarget, target.writeErr)

	r.On("RemoveAll", mock.Anything).Return(source.removeErr)

	return r
}

type __rwMock struct {
	mock.Mock
	rwcSource *__RWC
	rwcTarget *__RWC
}

func (rw *__rwMock) Write(patch string, flag int) (io.WriteCloser, error) {
	ret := rw.Called(patch, flag)

	return ret.Get(0).(io.WriteCloser), ret.Error(1)
}

func (rw *__rwMock) Read(patch string) (io.ReadCloser, error) {
	ret := rw.Called(patch)

	return ret.Get(0).(io.ReadCloser), ret.Error(1)
}

func (rw *__rwMock) RemoveAll(_ string) error {
	ret := rw.Called()

	return ret.Error(0)
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

// testWriter - write to buffer.
type __testWriter struct {
	buffer bytes.Buffer
}

// Write - write to buffer.
func (n *__testWriter) Write(p []byte) (int, error) {
	return n.buffer.Write(p)
}

// String - return buffer as string.
func (n *__testWriter) String() string {
	s := convert.B2S(n.buffer.Bytes())
	n.buffer.Reset()

	return s
}
