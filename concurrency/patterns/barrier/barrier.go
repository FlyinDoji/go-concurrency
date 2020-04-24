package main

import (
	"sync"
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

/*
 	Barrier is a generalized rendezvous that can sync parts of code between N threads
	reusableBarriers allow calling threads to sync at each Wait() call before proceeding further
	Uses two semaphores:
		First semaphore syncs the threads at the Wait() call
		Second semaphore returns the barrier in the initial state before threads can proceed

*/

type reusableBarrier struct {
	n             int
	threadCount   int
	firstBarrier  *semaphore
	secondBarrier *semaphore
	sync.Mutex
}

func newReusableBarrier(n int) *reusableBarrier {
	return &reusableBarrier{n: n, threadCount: 0, firstBarrier: newSemaphore(0), secondBarrier: newSemaphore(0)}
}

func (rb *reusableBarrier) first() {

	rb.Lock()
	rb.threadCount++
	if rb.threadCount == rb.n {
		for i := 0; i < rb.n; i++ {
			rb.firstBarrier.Signal()
		}
	}
	rb.Unlock()
	rb.firstBarrier.Wait()
}

func (rb *reusableBarrier) second() {

	rb.Lock()
	rb.threadCount--
	if rb.threadCount == 0 {
		for i := 0; i < rb.n; i++ {
			rb.secondBarrier.Signal()
		}
	}
	rb.Unlock()
	rb.secondBarrier.Wait()
}

func (rb *reusableBarrier) Wait() {
	// Both barriers locked

	// Block at 1st barrier, Nth thread unblocks this call for everyone
	rb.first()
	// Block at 2nd barrier, Nth thread locks first before unblocking second for everyone
	rb.second()
	// Nth thread to finish will lock second before unlocking first again.
}
