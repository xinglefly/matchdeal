package model

import (
	"time"
	"fmt"
	"math/rand"
)

type Maker struct {
	Price int
	Num   int
}

func createGorutingMaker() chan Maker {
	c := make(chan Maker)
	for i := 0; i <= 4; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				m := Maker{
					Price: rand.Intn(100),
					Num:   rand.Intn(6) + 1,
				}
				c <- m
			}
		}(i)
	}

	return c
}

func doWorkMaker(id int, c chan Maker) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Maker id %d receiver %d\n", id, n)
	}
}

func receiverMaker(id int) chan<- Maker {
	c := make(chan Maker)
	go doWorkMaker(id, c)
	return c
}

func CreateMaker() {
	m := createGorutingMaker()
	maker := receiverMaker(0)
	var queues []Maker
	var tike = time.Tick(time.Second)
	tm := time.After(10 * time.Second)

	for {
		var activeMaker chan<- Maker
		var makerValue Maker
		if len(queues) > 0 {
			activeMaker = maker
			makerValue = queues[0]
		}

		select {
		case n := <-m:
			queues = append(queues, n)
			fmt.Println("maker[]:", queues)
		case activeMaker <- makerValue:
			queues = queues[1:]
		case <-tike:
			fmt.Printf("makerQueues:%d \n", len(queues))
		case <-tm:
			fmt.Println("exit maker.")
			return
		}

	}

}
