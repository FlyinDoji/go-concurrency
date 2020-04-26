package main

import (
	"fmt"
	"time"
	"tutorials/concurrency/concurrency/patterns/semaphore"
)

//We want thread A to print "A2" after thread B has printed "B1"
//Similarly we want thread B to print "B2" after thread A has printed "A1"
//This situation is called a rendezvous, because the threads have 'meeting point'
//While "A1" and "B1" can be printed in any order, we want the threads to wait for each other
//before proceeding with the final instruction.

func A(a *semaphore.Semaphore, b *semaphore.Semaphore) {
	fmt.Println("A1")
	a.Signal() //Signal to B that A arrived at the rendezvous
	b.Wait()   //Block until B arrives
	fmt.Println("A2")
}
func B(a *semaphore.Semaphore, b *semaphore.Semaphore) {
	fmt.Println("B1")
	b.Signal() //Signal to A that B arrived at the rendezvous
	a.Wait()   //Block until A arrives
	fmt.Println("B2")
}

func main() {
	a := semaphore.NewSemaphore(0)
	b := semaphore.NewSemaphore(0)

	go A(a, b)
	go B(a, b)
	<-time.After(time.Second * 1)

}
