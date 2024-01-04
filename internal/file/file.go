package file

import (
	"io"
	"os"
)

// Reader is a file reader.
type Reader interface {
	Read(string) (io.ReadCloser, error)
}

// Writer is a file writer.
type Writer interface {
	Write(patch string, flag int) (io.WriteCloser, error)
}

// Remover is a file remover.
type Remover interface {
	Remove(string) error
}

// ReadWriter is a file reader and writer.
type ReadWriter interface {
	Reader
	Writer
}

// FS is a file system operations wrapper.
// Wrap for testing.
type FS struct {
	write  func(string, int) (io.WriteCloser, error)
	read   func(string) (io.ReadCloser, error)
	remove func(string) error
}

// FSOption is a file system option.
type FSOption struct {
	Write  func(string, int) (io.WriteCloser, error)
	Read   func(string) (io.ReadCloser, error)
	Remove func(string) error
}

// NewFS returns a new file system.
func NewFS(opts ...func(*FSOption)) *FS {
	fs := &FSOption{
		Write: func(p string, flag int) (io.WriteCloser, error) {
			return os.OpenFile(p, flag, 0o644)
		},
		Read: func(p string) (io.ReadCloser, error) {
			return os.OpenFile(p, os.O_RDONLY, 0o644)
		},
		Remove: os.Remove,
	}

	for _, opt := range opts {
		opt(fs)
	}

	return &FS{
		write:  fs.Write,
		read:   fs.Read,
		remove: fs.Remove,
	}
}

// Write returns a file for writing.
// flag:
//   - os.O_CREATE|os.O_WRONLY|os.O_TRUNC for rewrite file
//   - os.O_CREATE|os.O_APPEND|os.O_WRONLY for append file
func (f FS) Write(p string, flag int) (io.WriteCloser, error) {
	return f.write(p, flag)
}

// Read returns a file for reading.
func (f FS) Read(p string) (io.ReadCloser, error) {
	return f.read(p)
}

// Remove removes a file.
func (f FS) Remove(p string) error {
	return f.remove(p)
}
