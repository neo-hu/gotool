package sys

import "runtime"

const defaultUnixPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

func DefaultPathEnv(os string) string {
	if runtime.GOOS == "windows" {
		if os != runtime.GOOS {
			return defaultUnixPathEnv
		}
		// Deliberately empty on Windows containers on Windows as the default path will be set by
		// the container. Docker has no context of what the default path should be.
		return ""
	}
	return defaultUnixPathEnv

}