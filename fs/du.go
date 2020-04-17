package fs

type Usage struct {
	Inodes int64
	Size   int64
}

func DiskUsage(roots ...string) (Usage, error) {
	return diskUsage(roots...)
}
