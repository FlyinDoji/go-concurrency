package main

import (
	"fmt"
	"math/rand"
	"time"
	"tutorials/concurrency/concurrency/idiomatic"
)

type buffer []int
type eventBuffer struct {
	buffer
	c int
}

func (b *eventBuffer) Add(x int) {
	b.buffer[b.c] = x
	b.c++
}
func (b *eventBuffer) Get() int {
	b.c--
	return b.buffer[b.c]
}
func (b *eventBuffer) Empty() bool {
	return b.c == 0
}

func randomEventGenerator() int {
	r := rand.Intn(80) + 40
	time.Sleep(time.Millisecond * time.Duration(r))
	return r
}
func eventProcess(id, ev int) {
	for i := 0; i < 1000; i++ {
	}
	fmt.Println(id, "Processed", ev)
}

func producer(q *eventBuffer, space, items, mutex idiomatic.Semaphore) {
	for {
		ev := randomEventGenerator()
		space.Wait()
		mutex.Wait()
		q.Add(ev)
		mutex.Signal()
		items.Signal()
	}

}
func consumer(id int, q *eventBuffer, space, items, mutex idiomatic.Semaphore) {

	for {
		items.Wait()
		mutex.Wait()
		ev := q.Get()
		mutex.Signal()
		space.Signal()
		eventProcess(id, ev)
	}
}

func main() {

	size := 50

	eb := eventBuffer{buffer: make([]int, size), c: 0}

	items := idiomatic.NewSemaphore(0)
	space := idiomatic.NewSemaphore(size)
	mutex := idiomatic.NewSemaphore(1)

	go producer(&eb, space, items, mutex)
	go producer(&eb, space, items, mutex)
	go producer(&eb, space, items, mutex)
	go producer(&eb, space, items, mutex)

	go consumer(1, &eb, space, items, mutex)
	go consumer(2, &eb, space, items, mutex)

	<-time.After(time.Second * 2)

}
