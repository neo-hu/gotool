package file

import (
	"os"
	"syscall"
)

func SyncDir(dirName string) error {
	dir, err := os.OpenFile(dirName, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	err = dir.Sync()
	if pe, ok := err.(*os.PathError); ok && pe.Err == syscall.EINVAL {
		err = nil
	} else if err != nil {
		return err
	}

	return dir.Close()
}

// todo os rename 的重写
func RenameFile(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

