package model

import (
	"time"
	"fmt"
	"math/rand"
	"sort"
)

type Taker struct {
	Id    int
	Price int
	Num   int
}

type TakerWrapper struct {
	taker [] Taker
	by    func(p, q *Taker) bool
}

type TakerSort func(p, q *Taker) bool

func (t TakerWrapper) Len() int {
	return len(t.taker)
}

func (t TakerWrapper) Swap(i, j int) {
	t.taker[i], t.taker[j] = t.taker[j], t.taker[i]
}

func (t TakerWrapper) Less(i, j int) bool {
	return t.by(&t.taker[i], &t.taker[j])
}

func SortTaker(taker []Taker, by TakerSort) {
	sort.Sort(TakerWrapper{taker, by})
}

var tm = time.After(15 * time.Second)

func createGroutingTaker() chan Taker {
	c := make(chan Taker)
	for i := 0; i <= 4; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				t := Taker{
					Id:    ii,
					Price: rand.Intn(100),
					Num:   rand.Intn(3),
				}
				ii++
				c <- t
			}
		}(i)
	}
	return c
}

func doWorkTaker(id int, c chan Taker) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Taker id %d receiver id %d price %d num %d\n", id, n.Id, n.Price, n.Num)
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
			SortTaker(takerQueues, func(p, q *Taker) bool {
				return q.Price < p.Price
			})
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
