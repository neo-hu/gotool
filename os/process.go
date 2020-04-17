package os

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

type ErrProcessTimeout struct {
	path string
	pid  int
}

func (e *ErrProcessTimeout) Pid() int {
	return e.pid
}

func (e *ErrProcessTimeout) Error() string {
	return fmt.Sprintf("timeout, process:%s will be killed and %d not exit", e.path, e.pid)
}

// 启动一个进程，并且设定运行的时间，超时就kill -9
func RunWithTimeout(name string, timeout time.Duration, ctx context.Context, arg ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return "", err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		go func() {
			<-done
		}()
		// todo timeout kill group
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			return "", err
		}
		return "", &ErrProcessTimeout{cmd.Path, cmd.Process.Pid}
	case err = <-done:
		return out.String(), err
	}
}

// 判断进程是否存在
func Exists(pid int) bool {
	if err := syscall.Kill(pid, 0); err == nil {
		return true
	}
	return false
}

func Run(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	out.Reset()
	return out.String(), err
}