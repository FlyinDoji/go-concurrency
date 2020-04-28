package idiomatic

import (
	"sync"
	"testing"
	"time"
)

func TestMutexSemaphore(t *testing.T) {

	s := NewSemaphore(1)
	wg := sync.WaitGroup{}
	sharedCounter := 0
	iters := 2500
	n := 200

	testfun := func() {
		defer wg.Done()
		for j := 0; j < iters; j++ {
			s.Wait()
			sharedCounter++
			s.Signal()
		}

	}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go testfun()
	}

	wg.Wait()
	if sharedCounter != iters*n {
		t.Fatalf("Bad counter value:%d expected %d", sharedCounter, n*iters)
	}

}

func TestSemaphoreSignal(t *testing.T) {

	s := NewSemaphore(0)
	var done chan struct{}
	testfun := func() {
		defer close(done)
		s.Signal()
	}
	// goroutine should close the channel immediately
	done = make(chan struct{})
	go testfun()
	select {
	case <-done:
	case <-time.After(time.Second * 1):
		t.Fatalf("Signal() did not return immediately!")
	}

}

func TestSemaphoreWaitBeforeSignal(t *testing.T) {
	s := NewSemaphore(0)
	var done chan struct{}

	testfun := func() {
		defer close(done)
		s.Wait()
	}

	// goroutine should block
	done = make(chan struct{})
	go testfun()
	select {
	case <-done:
		t.Fatalf("Wait() did not block until signal!")
	case <-time.After(time.Second * 1):
	}
	s.Signal()

}

func TestSemaphoreWaitAfterSignal(t *testing.T) {
	s := NewSemaphore(0)
	var done chan struct{}

	s.Signal()

	testfun := func() {
		defer close(done)
		s.Wait()
	}
	// semaphore was signaled beforehand, goroutine should not block
	done = make(chan struct{})
	go testfun()
	select {
	case <-done:
	case <-time.After(time.Second * 1):
		t.Fatalf("Wait() blocked!")
	}

}

func TestBoundedSemaphoreWait(t *testing.T) {
	n := 100
	s := NewSemaphore(n)
	var done chan struct{}

	testfun := func() {
		s.Wait()
		done <- struct{}{}
	}
	// only n goroutines fill the channel
	done = make(chan struct{}, 2*n)
	// goroutines have time to finish work or block
	for i := 0; i < 2*n; i++ {
		go testfun()
		time.Sleep(time.Millisecond)
	}
	iter := 0
	ok := true
	for ok {
		select {
		case <-done:
			iter++
		default:
			ok = false
		}
	}

	if iter != n {
		t.Fatalf("Some threads have not blocked on Wait()! Limit: %d, total: %d", n, iter)
	}

}

func TestBarrier(t *testing.T) {

	n := 4
	iters := 100000
	b := NewBarrier(n)
	wg := sync.WaitGroup{}

	workload := make(chan struct{}, 2*n)

	testfun := func() {
		for i := 0; i < iters; i++ {
			// Each thread sends 1 item on the buffered channel
			workload <- struct{}{}
			// Sync here
			b.Wait()
			// All threads have finished sending before checking items in channel
			if len(workload) != n {
				t.Fatalf("Threads did not sync correctly! %d items in channel, expected %d", len(workload), n)
			}
			// Sync here
			b.Wait()
			// Each thread receives 1 item from the channel
			<-workload
			// Sync here
			b.Wait()
			// All threads have finished receiving before checking items in the channel
			if len(workload) != 0 {
				t.Fatalf("Threads did not sync correctly! %d items in channel, expected %d", len(workload), 0)
			}
			// Sync here
			b.Wait()
		}
		wg.Done()
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go testfun()
	}
	wg.Wait()
}
