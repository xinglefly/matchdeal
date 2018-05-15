package main

import (
	"fmt"
	"time"
)

type Maker int
type Taker int

func CreateMaker() chan int {
	c := make(chan int)
	for i := 0; i <= 20; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				c <- ii
			}
		}(i)
	}

	return c
}

func CreateTaker() chan int {
	c := make(chan int)
	for i := 0; i <= 10; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				c <- ii
			}
		}(i)
	}
	return c
}

//create Maker queue
func WorkerMaker(id int, c chan int) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Maker id %d receiver %d\n", id, n)
	}
}

func CreateMakerQueue(id int) chan<- int {
	c := make(chan int)
	go WorkerMaker(id, c)
	return c
}

//create taker queue
func WorkerTaker(id int, c chan int) {
	for n := range c {
		fmt.Printf("Taker id %d receiver %d\n", id, n)
	}
}

func CreateTakerQueue(id int) chan<- int {
	c := make(chan int)
	go WorkerTaker(id, c)
	return c
}

var m, t = CreateMaker(), CreateTaker()
var maker = CreateMakerQueue(0)
var makerQueues []int

var taker = CreateTakerQueue(0)
var takerQueues []int

func main() {

	tike := time.Tick(time.Second)
	tm := time.After(15 * time.Second)

	for {
		var activeMaker chan<- int
		var makerValue int
		if len(makerQueues) > 0 {
			activeMaker = maker
			makerValue = makerQueues[0]
		}

		var activeTaker chan<- int
		var takerValue int
		if len(takerQueues) > 0 {
			activeTaker = taker
			takerValue = takerQueues[0]
		}

		select {
		case n := <-m:
			makerQueues = append(makerQueues, n)
		case activeMaker <- makerValue:
			makerQueues = makerQueues[1:]
		case n := <-t:
			takerQueues = append(takerQueues, n)
		case activeTaker <- takerValue:
			takerQueues = takerQueues[1:]
		case <-tike:
			fmt.Printf("makerQueues:%d  takerQueues:%d \n:", len(makerQueues), len(takerQueues))
		case <-tm:
			fmt.Println("bye bye")
			return
		}

	}

}
