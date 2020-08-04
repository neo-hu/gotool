package wait

import (
	"sync"
)

type GroupWrapper struct {
	sync.WaitGroup
}

func (w *GroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
