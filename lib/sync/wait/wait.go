package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

func (w *Wait) Add(amount int) {
	w.wg.Add(amount)
}

func (w *Wait) Done() {
	w.wg.Done()
}

func (w *Wait) Wait() {
	w.wg.Wait()
}

func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}
