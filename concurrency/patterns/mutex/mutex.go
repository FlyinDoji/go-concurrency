// Little book of semaphores
// 1.5.3 Mutual exclusion with messages
// Like serialization, mutual exclusion can be implemented using message passing.
// For example, imagine that you and Bob operate a nuclear reactor that you
// monitor from remote stations. Most of the time, both of you are watching for
// warning lights, but you are both allowed to take a break for lunch. It doesn’t
// matter who eats lunch first, but it is very important that you don’t eat lunch
// at the same time, leaving the reactor unwatched!
// Puzzle: Figure out a system of message passing (phone calls) that enforces
// these restraints. Assume there are no clocks, and you cannot predict when lunch
// will start or how long it will last. What is the minimum number of messages
// that is required?

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const workdays = 5

func worker(name string, eatMe, eatOther *sync.Mutex, weekend *sync.WaitGroup) {

	defer weekend.Done()
	for i := 1; i <= workdays; i++ {
		eatMe.Lock()
		fmt.Println(name, "gone to lunch on day", i)
		time.Sleep(time.Second * time.Duration(rand.Intn(4)))
		fmt.Println(name, "finished eating.")
		eatOther.Unlock()
	}

}

func main() {

	kevinEats := sync.Mutex{}
	bobEats := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Each day one Kevin starts eating, after he is finished he signals Bob that he can now eat by unlocking eatOther
	// Kevin then cannot eat again until Bob has finished eating and unlocked eatOther
	// Order is important here, because they cannot start eating at the same time in Day 1, one worker needs to ensure he first waits for the other to finish
	// Hence the bobEats.Lock() that allows Kevin to eat first on day 1
	bobEats.Lock()
	go worker("bob", &bobEats, &kevinEats, &wg)
	go worker("kevin", &kevinEats, &bobEats, &wg)

	wg.Wait()
}
