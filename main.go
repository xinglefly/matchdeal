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
	orderId   = 1
)

//生成交易订单
func CreateOrder(op string, m model.Maker, t model.Taker) {
	//fmt.Println("op-->", op)
	if ok := model.MathQ(t.Price, m.Price); ok && m.Price <= t.Price {
		//TODO 数量相等时，出栈
		if m.Num == t.Num { //完全撮合
			PopQueues(op, m, t)
			PushOrder(op, t, m)
		} else if strings.EqualFold(op, "buy") && t.Num < m.Num {
			//TODO 买单数量 < 卖单数量，t 不入栈，makerQueue数量减少（先出栈再入栈 | ）
			PopQueues(op, m, t)
			//insertMakerQueue(model.UpdateMaker(m, m.Num - t.Num))
			fmt.Println("Update Maker-->", m, model.UpdateMaker(m, m.Num - t.Num))
			PushOrder(op, t, model.UpdateMaker(m, m.Num - t.Num))

		} else if strings.EqualFold(op, "buy") && t.Num > m.Num {
			//TODO 买单数量 > 卖单数量，m 出栈后 再 递归一次

		} else if strings.EqualFold(op, "sale") && t.Num > m.Num {
			//TODO 卖单 < 买单数量， m 不入栈， TakerQueue数量减少
			PopQueues(op, m, t)
			//insertTakerQueues(model.UpdateTaker(t, t.Num - m.Num))
			fmt.Println("Update Taker-->", t, model.UpdateTaker(t, t.Num - m.Num))

			PushOrder(op, model.UpdateTaker(t, t.Num - m.Num), m)

		} else if strings.EqualFold(op, "sale") && t.Num > m.Num {
			//TODO  卖单 > 买单， 递归一次
		}

	} else {
		if strings.EqualFold(op, "buy") {
			insertTakerQueues(t)
		} else {
			insertMakerQueue(m)
		}

	}
}
func PushOrder(op string, t model.Taker, m model.Maker) {
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

	model.SortTPrice2Time()
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
	model.SortMPrice2Time()
}

func main() {
	for {
		select {
		case n := <-dataTaker:
			MatchTaker(n)
		case activeTaker <- takerValue:
			model.QueuesTaker = model.QueuesTaker[1:]
			model.SortTPrice2Time()
		case n := <-dataMaker:
			MatchMaker(n)

		case activeMaker <- makerValue:
			model.QueuesMaker = model.QueuesMaker[1:]
			model.SortMPrice2Time()
		case <-tick:
			fmt.Printf("takerQueues:%d \n:", len(model.QueuesTaker))
			fmt.Printf("makerQueues:%d \n:", len(model.QueuesMaker))
		case <-tm:
			model.SortTPrice2Time()
			model.SortMPrice2Time()
			fmt.Println("taker[]", model.QueuesTaker)
			fmt.Println("maker[]:", model.QueuesMaker)
			fmt.Println("exit program!")
			return
		}
	}

}
