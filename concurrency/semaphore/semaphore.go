package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Semaphore implementation

type semaphore struct {
	capacity int

	count int
	sync.Mutex
	condition chan bool
}

// Acquire blocks the caller if thread capacity has been reached or returns otherwise
func (s *semaphore) Acquire() {
	s.Lock()
	defer s.Unlock()

	// Blocking the caller without keeping the Semaphore locked to other threads
	if s.count == s.capacity {
		s.Unlock()
		<-s.condition
		s.Lock()
	}

	s.count++

}

// Release uses a non-blocking send to prevent deadlock when there are no more threads waiting on the channel in Acquire
func (s *semaphore) Release() {
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

func worker(id int, sem *semaphore, wg *sync.WaitGroup) {
	sem.Acquire()
	n := time.Duration(rand.Intn(4)+1) * time.Second
	fmt.Println(id, "acquired sem for ", n)
	time.Sleep(n)
	sem.Release()
	wg.Done()

}

func main() {

	threads := 10
	sem := newSemaphore(3)
	wg := sync.WaitGroup{}

	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go worker(i, sem, &wg)
	}

	wg.Wait()
	threads = 25

	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go worker(i, sem, &wg)
	}
	wg.Wait()
}
