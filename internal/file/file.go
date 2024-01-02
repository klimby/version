package file

import (
	"io"
	"os"
)

type Reader interface {
	Read(string) (io.ReadCloser, error)
}

type Writer interface {
	Write(patch string, flag int) (io.WriteCloser, error)
}

type Remover interface {
	Remove(string) error
}

type ReadWriter interface {
	Reader
	Writer
}

type FS struct {
	write  func(string, int) (io.WriteCloser, error)
	read   func(string) (io.ReadCloser, error)
	remove func(string) error
}

type FSOption struct {
	Write  func(string, int) (io.WriteCloser, error)
	Read   func(string) (io.ReadCloser, error)
	Remove func(string) error
}

func NewFS(opts ...func(*FSOption)) *FS {
	fs := &FSOption{
		/*Write: func(p string, flag int) (io.WriteCloser, error) {
			return os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		},*/
		Write: func(p string, flag int) (io.WriteCloser, error) {
			return os.OpenFile(p, flag, 0644)
		},
		Read: func(p string) (io.ReadCloser, error) {
			return os.OpenFile(p, os.O_RDONLY, 0644)
		},
		/*Append: func(p string) (io.WriteCloser, error) {
			return os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		},*/
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
