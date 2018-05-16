package model

import (
	"time"
	"fmt"
	"math/rand"
)

type Taker struct {
	Price int
	Num   int
}

var tm = time.After(15 * time.Second)

func createGroutingTaker() chan Taker {
	c := make(chan Taker)
	for i := 0; i <= 4; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				t := Taker{
					Price: rand.Intn(100),
					Num:   rand.Intn(3),
				}
				c <- t
			}
		}(i)
	}
	return c
}

func doWorkTaker(id int, c chan Taker) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Taker id %d receiver price %d num %d\n", id, n.Price, n.Num)
	}
}

func receiverTaker(id int) chan<- Taker {
	c := make(chan Taker)
	go doWorkTaker(id, c)
	return c
}

func CreateTaker() {
	data := createGroutingTaker()
	taker := receiverTaker(0)
	var takerQueues []Taker

	tick := time.Tick(time.Second)

	for {
		var activeTaker chan<- Taker
		var takerValue Taker
		if len(takerQueues) > 0 {
			activeTaker = taker
			takerValue = takerQueues[0]
		}

		select {
		case n := <-data:
			takerQueues = append(takerQueues, n)
			fmt.Println("taker[]", takerQueues)
		case activeTaker <- takerValue:
			takerQueues = takerQueues[1:]
		case <-tick:
			fmt.Printf("takerQueues:%d \n:", len(takerQueues))
		case <-tm:
			fmt.Println("exit taker!")
			return
		}

	}
}
