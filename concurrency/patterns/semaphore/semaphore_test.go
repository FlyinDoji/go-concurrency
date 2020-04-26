package semaphore

import (
	"sync"
	"testing"
)

func TestMutexSemaphore(t *testing.T) {

	s := NewSemaphore(1)
	wg := sync.WaitGroup{}
	sharedCounter := 0
	iters := 2500
	n := 200

	testfun := func(mutex *Semaphore) {
		defer wg.Done()
		for j := 0; j < iters; j++ {
			s.Wait()
			sharedCounter++
			s.Signal()
		}

	}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go testfun(s)
	}

	wg.Wait()
	if sharedCounter != iters*n {
		t.Fatalf("Bad counter value:%d expected %d", sharedCounter, n*iters)
	}

}
