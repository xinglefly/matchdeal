package main

import (
	"time"
	"fmt"
	"matchdeal/model"
	"strings"
)

//TODO 后期重构用接口
type Goruting interface {
	CreateGoruting() chan int
	Receiver(id int) chan<- int
}

var (
	tm   = time.After(10 * time.Second)
	tick = time.Tick(time.Second)
	//Taker
	dataTaker = model.CreateGroutingTaker()
	taker     = model.ReceiverTaker(0)

	//Maker
	dataMaker = model.CreateGorutingMaker()
	maker     = model.ReceiverMaker(0)

	orderId = 1
)

//生成交易订单
func CreateOrder(op string, m model.Maker, t model.Taker) {
	//fmt.Println("op-->", op)
	if ok := model.MathQ(t.Price, m.Price); ok && m.Price <= t.Price {
		PopQueues(op, m, t)
		//币币交易规则限制
		push := model.Order{
			Id: orderId,
			Op: op,
			T:  t,
			M:  m,
		}
		//fmt.Println(push)
		model.OrderQueues = append(model.OrderQueues, push)
		orderId++
	} else {
		if strings.EqualFold(op, "buy") {
			insertTakerQueues(t)
		} else {
			insertMakerQueue(m)
		}

	}
}

var activeTaker chan<- model.Taker
var takerValue model.Taker

var activeMaker chan<- model.Maker
var makerValue model.Maker

//Pop Queues logic
func PopQueues(op string, m model.Maker, t model.Taker) {
	if strings.EqualFold(op, "buy") {
		activeMaker = maker
		makerValue = m
	} else {
		activeTaker = taker
		takerValue = t
	}
}

//买单 ——> 卖单队列中匹配
func MatchTaker(t model.Taker) {
	if len(model.QueuesMaker) > 0 {
		m := model.QueuesMaker[0]
		//fmt.Println("buy[]", m, t.Price)
		CreateOrder("buy", m, t)
	} else {
		insertTakerQueues(t)
	}
}

//插入到买单队列中
func insertTakerQueues(t model.Taker) {
	time.Sleep(15 * time.Millisecond)
	model.QueuesTaker = append(model.QueuesTaker, t)

	model.SortTaker(model.QueuesTaker, func(p, q *model.Taker) bool {
		return q.Price < p.Price
	})
	model.SortTaker(model.QueuesTaker, func(p, q *model.Taker) bool {
		if p.Price == q.Price {
			return p.Created < q.Created
		}
		return false
	})
}

//卖单 ——> 买单队列中匹配
func MatchMaker(m model.Maker) {
	if len(model.QueuesTaker) > 0 {
		t := model.QueuesTaker[0]
		CreateOrder("sale", m, t)
	} else {
		insertMakerQueue(m)
	}
}

func insertMakerQueue(m model.Maker) {
	time.Sleep(15 * time.Millisecond)
	model.QueuesMaker = append(model.QueuesMaker, m)
	model.SortMaker(model.QueuesMaker, func(q, p *model.Maker) bool {
		return p.Price < q.Price
	})

	model.SortMaker(model.QueuesMaker, func(q, p *model.Maker) bool {
		if p.Price == q.Price {
			return p.Created < q.Created
		}
		return false
	})
}

func main() {
	for {
		select {
		case n := <-dataTaker:
			MatchTaker(n)
		case activeTaker <- takerValue:
			model.QueuesTaker = model.QueuesTaker[1:]
			model.SortTaker(model.QueuesTaker, func(p, q *model.Taker) bool {
				return q.Price < p.Price
			})
		case n := <-dataMaker:
			MatchMaker(n)

		case activeMaker <- makerValue:
			model.QueuesMaker = model.QueuesMaker[1:]
			model.SortMaker(model.QueuesMaker, func(q, p *model.Maker) bool {
				return p.Price < q.Price
			})
		case <-tick:
			fmt.Printf("takerQueues:%d \n:", len(model.QueuesTaker))
			fmt.Printf("makerQueues:%d \n:", len(model.QueuesMaker))
		case <-tm:
			fmt.Println("taker[]", model.QueuesTaker)
			fmt.Println("maker[]:", model.QueuesMaker)
			fmt.Println("exit program!")
			return
		}
	}

}
