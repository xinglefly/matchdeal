package model

import (
	"time"
	"fmt"
)

type Maker int

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
		time.Sleep(2 * time.Second)
		fmt.Printf("Maker id %d receiver %d\n", id, n)
	}
}

func receiver(id int) chan<- int {
	c := make(chan int)
	go doWork(id, c)
	return c
}

func CreateMaker() {
	m := createGrouting()
	maker := receiver(0)
	var queues []int
	var tike = time.Tick(time.Second)

	for {
		var activeMaker chan<- int
		var makerValue int
		if len(queues) > 0 {
			activeMaker = maker
			makerValue = queues[0]
		}

		select {
		case n := <-m:
			queues = append(queues, n)
		case activeMaker <- makerValue:
			queues = queues[1:]
			fmt.Printf("makerQueues:%d \n", len(queues))
		case <-tike:
			//fmt.Printf("makerQueues:%d \n", len(makerQueues))
		}

	}

}
