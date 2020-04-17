package pool

import (
	"bytes"
	"testing"
)

func TestPool_Buffer(t *testing.T) {
	pool := NewPool(10, func(args interface{}) interface{} {
		return new(bytes.Buffer)
	})

	item := pool.Get(nil).(*bytes.Buffer)
	item.Write(make([]byte, 1024))
	//cap, len := item.Cap(), item.Len()
	pool.Put(item)

	item1 := pool.Get(nil).(*bytes.Buffer)
	item1.Reset()
	if item1.Cap() != item.Cap() {
		t.Fatalf("get item got %v, exp: %v", item.Cap(), item1.Cap())
	}
}
