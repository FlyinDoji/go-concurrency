package main

import (
	"sync"
	"tutorials/concurrency/concurrency/patterns/semaphore"
)

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
	firstBarrier  *semaphore.Semaphore
	secondBarrier *semaphore.Semaphore
	sync.Mutex
}

func newReusableBarrier(n int) *reusableBarrier {
	return &reusableBarrier{n: n, threadCount: 0, firstBarrier: semaphore.NewSemaphore(0), secondBarrier: semaphore.NewSemaphore(0)}
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
