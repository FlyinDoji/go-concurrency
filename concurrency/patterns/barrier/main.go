package main

import (
	"fmt"
	"time"
)

func worker(id int, rb *reusableBarrier) {

	for {
		fmt.Println(id, "Rendezvous")
		rb.Wait()
		fmt.Println(id, "Critical")
		rb.Wait()
	}

}

func main() {

	nThreads := 4

	b := newReusableBarrier(nThreads)
	for i := 1; i <= nThreads; i++ {
		go worker(i, b)
	}

	<-time.After(time.Second * 2)

}
