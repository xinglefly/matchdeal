package model

import (
	"time"
	"fmt"
	"math/rand"
	"sort"
)

type Maker struct {
	Id    int //TODO 内存中存储手动++，后期放到数据库自增长
	Price int
	Num   int
}

type MakerWrapper struct {
	maker []Maker
	by    func(q, p *Maker) bool
}

type MakerSort func(q, p *Maker) bool

func (m MakerWrapper) Len() int {
	return len(m.maker)
}

func (m MakerWrapper) Swap(i, j int) {
	m.maker[i], m.maker[j] = m.maker[j], m.maker[i]
}

func (m MakerWrapper) Less(i, j int) bool {
	return m.by(&m.maker[i], &m.maker[j])
}

func SortMaker(maker []Maker, by MakerSort) {
	sort.Sort(MakerWrapper{maker, by})
}

func CreateGorutingMaker() chan Maker {
	c := make(chan Maker)
	for i := 1; i < 5; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				m := Maker{
					Id:    ii,
					Price: rand.Intn(100),
					Num:   rand.Intn(6) + 1,
				}
				i++
				c <- m
			}
		}(i)
	}

	return c
}

func doWorkMaker(id int, c chan Maker) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Maker id %d receiver price %d num%d\n", n.Id, n.Price, n.Num)
	}
}

func ReceiverMaker(id int) chan<- Maker {
	c := make(chan Maker)
	go doWorkMaker(id, c)
	return c
}

func CreateMaker() {
	/*m := createGorutingMaker()
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
			SortMaker(queues, func(q, p *Maker) bool {
				return p.Price < q.Price
			})
			fmt.Println("maker[]:", queues)
		case activeMaker <- makerValue:
			queues = queues[1:]
		case <-tike:
			fmt.Printf("makerQueues:%d \n", len(queues))
		case <-tm:
			fmt.Println("exit maker.")
			return
		}

	}*/

}
