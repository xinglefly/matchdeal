package model

import (
	"time"
	"fmt"
	"math/rand"
	"sort"
	"matchdeal/common"
)

type Taker struct {
	Id      int
	Price   int
	Num     int
	Created string
}

var QueuesTaker []Taker

func (t Taker) String() string {
	return fmt.Sprintf("T[id=%d, p=%d,n=%d,t=%s]", t.Id, t.Price, t.Num, t.Created)
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

func SortPrice2Time(){

}


func CreateGroutingTaker() chan Taker {
	c := make(chan Taker)
	for i := 1; i < 40; i++ {
		go func(ii int) {
			for {
				time.Sleep(time.Duration(1500) * time.Millisecond)
				t := Taker{
					Id:      ii,
					Price:   rand.Intn(100),
					Num:     rand.Intn(3) + 1,
					Created: common.FormatTime(),
				}
				i++
				c <- t
			}
		}(i)
	}
	return c
}

func doWorkTaker(id int, c chan Taker) {
	for n := range c {
		time.Sleep(2 * time.Second)
		fmt.Printf("Taker id %d receiver price %d num %d\n", n.Id, n.Price, n.Num)
		fmt.Println("pop takerQueue[0]", n)
		fmt.Println("order[]", OrderQueues)
	}
}

func ReceiverTaker(id int) chan<- Taker {
	c := make(chan Taker)
	go doWorkTaker(id, c)
	return c
}
