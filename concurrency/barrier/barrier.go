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

func worker(id, threadCount int, waitCount *int, mutex *semaphore, barrier, barrier2 *semaphore) {

	for {
		println(id, "rendezvous")
		mutex.Wait()
		*waitCount--
		// Nth thread unlocks the first barrier and locks the 2nd
		if *waitCount == 0 {
			barrier2.Wait()
			barrier.Signal()
		}
		mutex.Signal()

		// Block here until everyone reached rendezvous and thread N signals the barrier
		barrier.Wait()
		// Thread N-i unlocks the barrier by decrementing the semaphore to -1 and waking up the N-i-1 thread.
		// Thread N-i-1 will increment the barrier to 0 and signal the N-2 thread, until all threads pass.
		// Note: Thread N will call Signal() two times, this will cause the semaphore to be decremented one extra
		// The barrier will not reset after all threads have executed the critical part
		// This potentially allows a thread to loop around several times, effectively getting ahead of the others
		barrier.Signal()
		println(id, "critical after")

		// Barrier reset logic here
		// Requires a second barrier that blocks all threads until the first barrier is reset
		mutex.Wait()
		*waitCount++
		// Nth thread locks the first barrier and unlocks the 2nd
		if *waitCount == threadCount {
			barrier.Wait()
			barrier2.Signal()
		}
		mutex.Signal()

		//Block here again until barrier 2 is unlocked
		barrier2.Wait()
		barrier2.Signal()
	}
}

func main() {

	nThreads := 4
	waitCount := nThreads
	mutex := newSemaphore(1)
	barrier := newSemaphore(0)
	barrier2 := newSemaphore(1)

	for i := 1; i <= nThreads; i++ {
		go worker(i, waitCount, &waitCount, mutex, barrier, barrier2)
	}

	<-time.After(time.Second * 5)

}
