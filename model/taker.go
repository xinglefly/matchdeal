package model

import (
	"time"
	"fmt"
)

type Taker int

var tm = time.After(15 * time.Second)

func createGrouting() chan int {
	c := make(chan int)
	for i := 0; i <= 4; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				c <- ii
			}
		}(i)
	}
	return c
}

func doWork(id int, c chan int) {
	for n := range c {
		fmt.Printf("Taker id %d receiver %d\n", id, n)
	}
}

func receiver(id int) chan<- int {
	c := make(chan int)
	go doWork(id, c)
	return c
}

func CreateTaker() {
	data := createGrouting()
	taker := receiver(0)
	var queues []int

	tick := time.Tick(time.Second)
	for {
		var activeTaker chan<- int
		var takerValue int
		if len(queues) > 0 {
			activeTaker = taker
			takerValue = queues[0]
		}

		select {
		case n := <-data:
			queues = append(queues, n)
		case activeTaker <- takerValue:
			queues = queues[1:]
		case <-tick:
			fmt.Printf("takerQueues:%d \n:", len(queues))
		case <-tm:
			fmt.Println("testing exit project!")
			return
		}

	}
}
