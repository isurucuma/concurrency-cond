package main

import (
	"fmt"
	"sync"
)

//func main() {
//	money := 0
//	mutex := sync.Mutex{}
//	cond := sync.NewCond(&mutex)
//	go stingy(&money, cond)
//	go spendy(&money, cond)
//	time.Sleep(2 * time.Second)
//	mutex.Lock()
//	fmt.Println("Money: ", money)
//	mutex.Unlock()
//}
//
//func stingy(money *int, cond *sync.Cond) {
//	for i := 0; i < 10000000; i++ {
//		cond.L.Lock()
//		*money += 1
//		cond.Signal()
//		cond.L.Unlock()
//	}
//	fmt.Println("Stingy Done")
//}
//
//func spendy(money *int, cond *sync.Cond) {
//	for i := 0; i < 200000; i++ {
//		cond.L.Lock()
//		for *money < 50 {
//			cond.Wait()
//		}
//		*money -= 50
//		if *money < 0 {
//			fmt.Println("Money is negative!")
//			os.Exit(1)
//		}
//		cond.L.Unlock()
//	}
//	fmt.Println("Spendy Done")
//}

// in this we have used the mutex to lock the condition variable before Signalling so that after the wait() hit only the signal() will be executed
//func doWork(cond *sync.Cond) {
//	fmt.Println("Work started")
//	fmt.Println("Work done")
//	cond.L.Lock()
//	cond.Signal()
//	cond.L.Unlock()
//}
//
//func main() {
//	m := sync.Mutex{}
//	cond := sync.NewCond(&m)
//	cond.L.Lock()
//	for i := 0; i < 50000; i++ {
//		go doWork(cond)
//		fmt.Println("Waiting for child work to finish")
//		cond.Wait()
//		fmt.Println("Child work finished")
//	}
//	cond.L.Unlock()
//	fmt.Println("All work done")
//}

// here in this example we are using the condition variable to wait for all players to be ready before starting the game
// here use of Broadcast is important as it will cause all the waiting goroutines to wake up, if we use signal then only
// one goroutine will wake up and the others will still be waiting
//func main() {
//	playersRemaining := 4
//	cond := sync.NewCond(&sync.Mutex{})
//	for i := 0; i < 4; i++ {
//		go playerHandler(cond, &playersRemaining, i)
//		time.Sleep(1 * time.Second)
//	}
//	fmt.Println("All players are ready")
//}
//
//func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
//	cond.L.Lock()
//	fmt.Println("Player", playerId, "is ready")
//	*playersRemaining--
//	if *playersRemaining == 0 {
//		cond.Broadcast()
//	}
//	if *playersRemaining > 0 {
//		cond.Wait()
//	}
//	cond.L.Unlock()
//}

type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{
		permits: permits,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}

func main() {
	sem := NewSemaphore(0)
	for i := 0; i < 5; i++ {
		go doWork(sem)
		fmt.Println("Waiting for work to be done")
		sem.Acquire()
		fmt.Println("Child work done")
	}
}

func doWork(sem *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work done")
	sem.Release()
}
