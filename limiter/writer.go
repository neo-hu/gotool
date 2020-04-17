package limiter

import (
	"context"
	"io"
)

type Writer interface {
	io.Writer
	io.Closer
	Name() string
	Sync() error
}

type rateWriter struct {
	w       io.WriteCloser
	limiter Rate
	ctx     context.Context
}

// 限制速率的write， bytesPerSec: 每秒最多写入bytesPerSec，如果bytesPerSec <=0 不限制速率
func NewWriter(w io.WriteCloser, bytesPerSec, burstLimit int) Writer {
	var limiter Rate
	if bytesPerSec > 0 {
		limiter = NewRate(bytesPerSec, burstLimit)
	}
	return &rateWriter{
		limiter: limiter,
		w:       w,
		ctx:     context.Background(),
	}
}

func (s *rateWriter) Sync() error {
	if f, ok := s.w.(interface {
		Sync() error
	}); ok {
		return f.Sync()
	}
	return nil
}
func (s *rateWriter) Close() error {
	return s.w.Close()
}
func (s *rateWriter) Name() string {
	if f, ok := s.w.(interface {
		Name() string
	}); ok {
		return f.Name()
	}
	return ""
}

func (s *rateWriter) Write(b []byte) (int, error) {
	if s.limiter == nil {
		return s.w.Write(b)
	}

	var n int
	for n < len(b) {
		wantToWriteN := len(b[n:])
		if wantToWriteN > s.limiter.Burst() {
			wantToWriteN = s.limiter.Burst()
		}
		wroteN, err := s.w.Write(b[n : n+wantToWriteN])
		if err != nil {
			return n, err
		}
		n += wroteN
		if err := s.limiter.WaitN(s.ctx, wroteN); err != nil {
			return n, err
		}
	}
	return n, nil
}
