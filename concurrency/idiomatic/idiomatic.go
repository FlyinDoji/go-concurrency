package idiomatic

// Go primitives allow an idiomatic implementation of the semaphore

// Semaphore wraps an empty struct channel
type Semaphore chan struct{}

func (s Semaphore) wake() {
	<-s
}

// Wait return or block until a Signal is sent
func (s Semaphore) Wait() {
	s <- struct{}{}
}

// Signal will schedule a new goroutine to wake one goroutine that blocks on Wait()
func (s Semaphore) Signal() {
	go s.wake()
}

func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}

// Barrier is a generalized rendezvous, used to sync N threads at certain execution points
type Barrier struct {
	n int

	count int
	mutex Semaphore
	in    Semaphore
	out   Semaphore
}

func (b *Barrier) enter() {

	b.mutex.Wait()
	b.count++
	if b.count == b.n {
		for i := 0; i < b.n; i++ {
			b.in.Signal()
		}
	}
	b.mutex.Signal()

	b.in.Wait()

}

func (b *Barrier) exit() {

	b.mutex.Wait()
	b.count--

	if b.count == 0 {
		for i := 0; i < b.n; i++ {
			b.out.Signal()
		}
	}
	b.mutex.Signal()

	b.out.Wait()

}

func (b *Barrier) Wait() {
	b.enter()
	b.exit()
}

func NewBarrier(n int) *Barrier {
	return &Barrier{n: n, count: 0, mutex: NewSemaphore(1), in: NewSemaphore(0), out: NewSemaphore(0)}
}
