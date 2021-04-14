package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(ExpandTildeIgnoreErr(path))
	return err == nil || os.IsExist(err)
}

func ExpandTildeIgnoreErr(path string) string {
	np, err := ExpandTilde(path)
	if err != nil {
		return path
	}
	return np
}

// ~/ -> /home/
func ExpandTilde(path string) (string, error) {
	if path == "~" {
		return HomeDir()
	}

	path = filepath.FromSlash(path)
	if !strings.HasPrefix(path, fmt.Sprintf("~%c", os.PathSeparator)) {
		return path, nil
	}

	home, err := HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path[2:]), nil
}

func HomeDir() (string, error) {
	if runtime.GOOS == "windows" {
		home := filepath.Join(os.Getenv("HomeDrive"), os.Getenv("HomePath"))
		if home != "" {
			return home, nil
		}
	}
	return os.UserHomeDir()
}

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func TempFile(dir, prefix string) (f *os.File, err error) {
	return ioutil.TempFile(dir, prefix)
}

func CreateFile(path string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(ExpandTildeIgnoreErr(path), data, perm)
}

func IsDir(path string) bool {
	f, e := os.Stat(ExpandTildeIgnoreErr(path))
	if e != nil {
		return false
	}
	return f.IsDir()
}

func IsFile(path string) bool {
	f, e := os.Stat(ExpandTildeIgnoreErr(path))
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func SearchFile(filename string, paths ...string) (fullPath string, err error) {
	if path.IsAbs(filename) {
		fullPath = filename
		return
	}
	lastPath := ""
	for _, path1 := range paths {
		if fullPath = filepath.Join(path1, filename); IsExist(fullPath) {
			return
		}
		if lastPath == "" {
			lastPath = fullPath
		}
	}
	err = fmt.Errorf("%s not found in paths", strings.Join(paths, ","))
	return
}

func DirNum(dirPath string) (int, error) {
	if !IsExist(dirPath) {
		return 0, os.ErrNotExist
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return 0, err
	}
	return len(fs), nil
}
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

func RemoveDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, f := range files {
		err = os.RemoveAll(f)
		if err != nil {
			return err
		}
	}
	return os.RemoveAll(dir)
}
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}
