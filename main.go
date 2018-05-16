package main

import (
	"time"
	"fmt"
	"matchdeal/model"
)

//TODO 后期重构用接口
type Goruting interface {
	CreateGoruting() chan int
	Receiver(id int) chan<- int
}

func getTaker(taker model.Taker) {

}

var (
	tm   = time.After(15 * time.Second)
	tick = time.Tick(time.Second)
	//Taker
	dataTaker   = model.CreateGroutingTaker()
	taker       = model.ReceiverTaker(0)
	queuesTaker []model.Taker
	//Maker
	dataMaker   = model.CreateGorutingMaker()
	maker       = model.ReceiverMaker(0)
	queuesMaker []model.Maker
)

func main() {
	for {
		var activeTaker chan<- model.Taker
		var takerValue model.Taker
		if len(queuesTaker) > 0 {
			activeTaker = taker
			takerValue = queuesTaker[0]
		}

		var activeMaker chan<- model.Maker
		var makerValue model.Maker
		if len(queuesMaker) > 0 {
			activeMaker = maker
			makerValue = queuesMaker[0]
		}

		select {
		case n := <-dataTaker:
			queuesTaker = append(queuesTaker, n)
			model.SortTaker(queuesTaker, func(p, q *model.Taker) bool {
				return q.Price < p.Price
			})
		case activeTaker <- takerValue:
			queuesTaker = queuesTaker[1:]

		case n := <-dataMaker:
			queuesMaker = append(queuesMaker, n)
			model.SortMaker(queuesMaker, func(q, p *model.Maker) bool {
				return p.Price < q.Price
			})

		case activeMaker <- makerValue:
			queuesMaker = queuesMaker[1:]
		case <-tick:
			fmt.Printf("takerQueues:%d \n:", len(queuesTaker))
			fmt.Printf("makerQueues:%d \n:", len(queuesMaker))
		case <-tm:
			fmt.Println("taker[]", queuesTaker)
			fmt.Println("maker[]:", queuesMaker)
			fmt.Println("exit program!")
			return
		}
	}

}
