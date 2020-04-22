package main

import (
	"fmt"
	"sync"
	"time"
)

// The semaphore is a generalized lock on which the Acquire() method can be called by N threads
// consecutively before it blocks the caller until a Release() call.
// If initialized with a capacity of 0 it can be used as a signaling mechanism,
// with Acquire/Wait blocking until a Release/Signal.

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

//We want thread A to print "A2" after thread B has printed "B1"
//Similarly we want thread B to print "B2" after thread A has printed "A1"
//This situation is called a rendezvous, because the threads have 'meeting point'
//While "A1" and "B1" can be printed in any order, we want the threads to wait for each other
//before proceeding with the final instruction.

func A(a *semaphore, b *semaphore) {
	fmt.Println("A1")
	a.Signal() //Signal to B that A arrived at the rendezvous
	b.Wait()   //Block until B arrives
	fmt.Println("A2")
}
func B(a *semaphore, b *semaphore) {
	fmt.Println("B1")
	b.Signal() //Signal to A that B arrived at the rendezvous
	a.Wait()   //Block until A arrives
	fmt.Println("B2")
}

func main() {

	a := newSemaphore(0)
	b := newSemaphore(0)
	go A(a, b)
	go B(a, b)
	<-time.After(time.Second * 1)

}
