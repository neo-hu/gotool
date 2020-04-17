package limiter

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

func TestWriter_Limited(t *testing.T) {
	r := bytes.NewReader(bytes.Repeat([]byte{0}, 512*5))

	limit := 512
	w := NewWriter(nopWriteCloser{ioutil.Discard}, limit, 10*1024)

	start := time.Now()
	n, err := io.Copy(w, r)
	elapsed := time.Since(start)
	if err != nil {
		t.Error("copy error: ", err)
	}
	rate := float64(n) / elapsed.Seconds()
	if rate > float64(limit) {
		t.Errorf("rate limit mismath: exp %f, got %f", float64(limit), rate)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (d nopWriteCloser) Close() error { return nil }
