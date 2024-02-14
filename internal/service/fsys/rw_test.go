package fsys

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFS(t *testing.T) {
	// Create a temporary file.
	f, err := os.CreateTemp("", "sample")
	if err != nil {
		t.Fatal(err)
	}

	defer assert.NoError(t, f.Close(), "close file")

	fs := New(func(option *Option) {})

	str := "Hello, World!"

	if err := __FSWrite(fs, f, str); err != nil {
		t.Fatal(err)
	}

	// Read from the file.
	res, err := __FSRead(fs, f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, res, "read from file")

	assert.True(t, fs.Exists(f.Name()), "file exists")

	// Remove the file.
	assert.NoError(t, fs.RemoveAll(f.Name()), "remove file")
	// Twice!
	assert.NoError(t, fs.RemoveAll(f.Name()), "remove file")
}

func __FSWrite(fs *FS, f *os.File, data string) (err error) {
	// Write to the file.
	w, err := fs.Write(f.Name(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}

	defer func() {
		if e := w.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	// write to the file.
	if _, err = w.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func __FSRead(fs *FS, f *os.File) (_ string, err error) {
	// Read from the file
	r, err := fs.Read(f.Name())
	if err != nil {
		return "", err
	}

	defer func() {
		if e := r.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	// Read from the file.
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
