package fsys

import (
	"io"
	"os"
)

// FS is a file system operations wrapper.
// Wrap for testing.
type FS struct {
	write  func(string, int) (io.WriteCloser, error)
	read   func(string) (io.ReadCloser, error)
	remove func(string) error
	exists func(string) bool
}

// Option is a file system option.
type Option struct {
	Write  func(string, int) (io.WriteCloser, error)
	Read   func(string) (io.ReadCloser, error)
	Remove func(string) error
	Exists func(string) bool
}

// New returns a new file system.
func New(opts ...func(*Option)) *FS {
	o := &Option{
		Write: func(p string, flag int) (io.WriteCloser, error) {
			return os.OpenFile(p, flag, 0o644)
		},
		Read: func(p string) (io.ReadCloser, error) {
			return os.OpenFile(p, os.O_RDONLY, 0o644)
		},
		Remove: os.RemoveAll,
		Exists: func(s string) bool {
			_, err := os.Stat(s)

			return err == nil
		},
	}

	for _, opt := range opts {
		opt(o)
	}

	return &FS{
		write:  o.Write,
		read:   o.Read,
		remove: o.Remove,
		exists: o.Exists,
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

// RemoveAll removes a file.
func (f FS) RemoveAll(p string) error {
	return f.remove(p)
}

// Exists returns true if the file exists.
func (f FS) Exists(p string) bool {
	return f.exists(p)
}
