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

	orderId = 1
)

//生成交易订单
func CreateOrder(op string, m model.Maker, t model.Taker) {
	//fmt.Println("op-->", op)
	if ok := model.MathQ(t.Price, m.Price); ok && m.Price <= t.Price {
		PopQueues(op,m,t)
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
	} else{
		activeTaker = taker
		takerValue = t
	}
}

//买单 ——> 卖单队列中匹配
func MatchTaker(t model.Taker) {
	if len(queuesMaker) > 0 {
		m := queuesMaker[0]
		//fmt.Println("buy[]", m, t.Price)
		CreateOrder("buy", m, t)
	} else {
		insertTakerQueues(t)
	}
}

//插入到买单队列中
func insertTakerQueues(t model.Taker) {
	time.Sleep(15 * time.Millisecond)
	queuesTaker = append(queuesTaker, t)

	model.SortTaker(queuesTaker, func(p, q *model.Taker) bool {
		return q.Price < p.Price
	})
}

//卖单 ——> 买单队列中匹配
func MatchMaker(m model.Maker) {
	if len(queuesTaker) > 0 {
		t := queuesTaker[0]
		CreateOrder("sale", m, t)
	} else {
		insertMakerQueue(m)
	}
}

func insertMakerQueue(m model.Maker) {
	time.Sleep(15 * time.Millisecond)
	queuesMaker = append(queuesMaker, m)
	model.SortMaker(queuesMaker, func(q, p *model.Maker) bool {
		return p.Price < q.Price
	})
}




func main() {
	for {
		select {
		case n := <-dataTaker:
			MatchTaker(n)
		case activeTaker <- takerValue:
			queuesTaker = queuesTaker[1:]
			model.SortTaker(queuesTaker, func(p, q *model.Taker) bool {
				return q.Price < p.Price
			})
		case n := <-dataMaker:
			MatchMaker(n)

		case activeMaker <- makerValue:
			queuesMaker = queuesMaker[1:]
			model.SortMaker(queuesMaker, func(q, p *model.Maker) bool {
				return p.Price < q.Price
			})
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
