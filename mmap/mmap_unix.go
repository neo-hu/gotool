package mmap

import (
	"os"
	"syscall"
)

// 将一个文件或者其它对象映射进内存
func Map(path string, sz int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	} else if fi.Size() == 0 {
		return nil, nil
	}
	if sz == 0 {
		sz = fi.Size()
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, int(sz), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return data, nil
}
