package main

import (
	"fmt"
	"math/rand"
	"time"
	"tutorials/concurrency/concurrency/idiomatic"
)

// Similar to producer-consumer,
// Read threads can access the critical section simultaneously
// Write threads require exclusive access to the critical section
// Readers exclude Writers
// Writers exclude Readers and Writers

func doStuff() {
	time.Sleep(time.Second * time.Duration((rand.Intn(1) + 1)))
}

func doCriticalStuff(msg string) {

	fmt.Println("Critical stuff - ", msg)
	doStuff()
}

func reader(m, empty, rt idiomatic.Semaphore, readers *int) {
	for {

		doStuff()
		rt.Wait()   // Block until a writer has finished
		rt.Signal() // Move through after writer has unblocked
		m.Wait()
		*readers++
		if *readers == 1 {
			empty.Wait()
		}
		m.Signal()

		doCriticalStuff("reader")

		m.Wait()
		*readers--
		if *readers == 0 {
			empty.Signal()
		}
		m.Signal()
	}
}
func writer(empty, rt idiomatic.Semaphore) {
	for {
		doStuff()
		// Without this section, writers can starve as readers cycle through without ever completely leaving the protected section
		//
		rt.Wait()    // Lock turnstile even if the readers haven't all finished
		empty.Wait() // Wait for readers to exit
		//
		doCriticalStuff("writer")
		empty.Signal()
		rt.Signal()
	}
}

func main() {
	readers := new(int)
	turnstile := idiomatic.NewSemaphore(1)
	mutex := idiomatic.NewSemaphore(1)
	empty := idiomatic.NewSemaphore(1)

	go reader(mutex, empty, turnstile, readers)
	go reader(mutex, empty, turnstile, readers)
	go reader(mutex, empty, turnstile, readers)
	go reader(mutex, empty, turnstile, readers)
	go reader(mutex, empty, turnstile, readers)

	go writer(empty, turnstile)
	go writer(empty, turnstile)
	go writer(empty, turnstile)

	<-time.After(time.Second * 20)

}
