package main

import (
	"fmt"
	"sync"
)

// Can embed the mutex and lock the structure itself
type sharedCounter struct {
	count int
	sync.Mutex
}

func (sc *sharedCounter) IncrementAndGet() int {
	sc.Lock()
	defer sc.Unlock()
	sc.count++
	return sc.count
}
func (sc *sharedCounter) Get() int {
	sc.Lock()
	defer sc.Unlock()
	return sc.count
}

// Synchronization structures must always be passed by pointer
func worker(id int, sc *sharedCounter, wg *sync.WaitGroup) {

	fmt.Println(id, "incremented to: ", sc.IncrementAndGet())
	wg.Done()

}

func main() {

	var sc sharedCounter
	var wg sync.WaitGroup
	threadCount := 1000

	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go worker(i, &sc, &wg)
	}

	wg.Wait()
	fmt.Println("Main exit:", sc.Get())

}
