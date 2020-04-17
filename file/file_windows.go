package file

import "os"

func SyncDir(dirName string) error {
	return nil
}

// 判断是否存在，如果存在先删除
func RenameFile(oldpath, newpath string) error {
	if IsExist(newpath) {
		if err = os.Remove(newpath); nil != err {
			return err
		}
	}
	return os.Rename(oldpath, newpath)
}
