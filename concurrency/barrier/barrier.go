package main

import (
	"sync"
	"time"
)

type semaphore struct {
	capacity int

	count int
	sync.Mutex
	condition chan bool
}

func (s *semaphore) Wait() {
	s.Lock()
	defer s.Unlock()
	if s.count == s.capacity {
		s.Unlock()
		<-s.condition
		s.Lock()
	}

	s.count++

}
func (s *semaphore) Signal() {
	s.Lock()
	defer s.Unlock()
	s.count--
	select {
	case s.condition <- true:
	default:
	}

}

func newSemaphore(capacity int) *semaphore {
	return &semaphore{count: 0, capacity: capacity, condition: make(chan bool)}
}

// Barrier is a generalized rendezvous with N threads
// The part before rendezvous can be executed in any order
// No thread may proceed to the critical part until all threads have reached the rendezvous

func worker(id int, waitCount *int, mutex *semaphore, barrier *semaphore) {

	println(id, "rendezvous")
	mutex.Wait()
	*waitCount--
	// Nth thread to finish can signal the barrier and wake the n-1 thread
	if *waitCount == 0 {
		barrier.Signal()
	}
	mutex.Signal()

	// Block here until everyone reached rendezvous and thread N signals the barrier
	barrier.Wait()
	// Thread N-i unlocks the barrier by decrementing the semaphore to -1 and waking up the N-i-1 thread.
	// Thread N-i-1 will increment the barrier to 0 and signal the N-2 thread, until all threads pass.
	// Note: Thread N will call Signal() two times, this will cause the semaphore to be decremented one extra
	// The barrier will not reset after all threads have executed the critical part
	barrier.Signal()
	println(id, "critical after")

	// Barrier reset logic here

}

func main() {

	nThreads := 5
	waitCount := nThreads
	mutex := newSemaphore(1)
	barrier := newSemaphore(0)

	for i := 1; i <= nThreads; i++ {
		go worker(i, &waitCount, mutex, barrier)
	}

	<-time.After(time.Second * 3)

	//The barrier in it's current form is not reusable, after all threads finish executing the counter is -1
	println(barrier.count)
}
