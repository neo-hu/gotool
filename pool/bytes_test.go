package pool

import "testing"

func TestLimitedBytePool_Put_MaxSize(t *testing.T) {
	exp := 0
	bp := NewLimitedBytes(1, exp)
	bp.Put(make([]byte, 1024)) // put一个大于maxSize的byte

	if got := cap(bp.Get(10)); got != exp {
		t.Fatalf("max cap size exceeded: got %v, exp %v", got, exp)
	}
}
