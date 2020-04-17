package file

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrClosed = errors.New("write to closed writer")
)

type atomicFileWriter struct {
	f    *os.File
	fn   string
	err  error
	perm os.FileMode
}

func NewAtomicFileWriter(filename string, perm os.FileMode) (io.WriteCloser, error) {
	f, err := ioutil.TempFile(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	abspath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	return &atomicFileWriter{
		f:    f,
		fn:   abspath,
		perm: perm,
	}, nil
}

func (w *atomicFileWriter) Write(dt []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.f.Write(dt)
	if err != nil {
		w.err = err
		w.f.Close()
	}
	return n, err
}

func (w *atomicFileWriter) Close() error {
	if w.err != nil {
		return w.err
	}
	if err := w.f.Sync(); err != nil {
		w.f.Close()
		return err
	}
	defer os.Remove(w.f.Name())
	if err := w.f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(w.f.Name(), w.perm); err != nil {
		return err
	}
	err := os.Rename(w.f.Name(), w.fn)
	if runtime.GOOS == "windows" && os.IsPermission(err) {
		_ = os.Chmod(w.fn, 0644)
		err = os.Rename(w.f.Name(), w.fn)
	}
	if err != nil {
		w.err = err
		return err
	}
	// fsync the directory too
	SyncDir(filepath.Dir(w.f.Name()))
	w.err = ErrClosed
	return nil
}
