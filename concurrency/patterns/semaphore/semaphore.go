package semaphore

import (
	"sync"
)

// Semaphore naive implementation
type Semaphore struct {
	capacity int

	count int
	sync.Mutex
	condition chan bool
}

// Wait returns immediately or blocks until a Signal(), increments the semaphore
func (s *Semaphore) Wait() {
	s.Lock()
	defer s.Unlock()

	// When waking, check that the condition still holds
	for s.count == s.capacity {
		s.Unlock()
		<-s.condition
		s.Lock()
	}

	s.count++

}

// Signal returns immediately, decrements the semaphore and wakes one waiting thread
func (s *Semaphore) Signal() {
	s.Lock()
	defer s.Unlock()

	s.count--

	select {
	case s.condition <- true:
	default:
	}

}

// NewSemaphore constructor
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{count: 0, capacity: n, condition: make(chan bool)}
}
