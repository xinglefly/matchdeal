package model

import (
	"strconv"
	"math"
	"fmt"
	"sync"
)

type Order struct {
	Id   int
	Op   string
	T    Taker
	M    Maker
	Lock sync.Mutex
}

var orderId = 1
var OrderQueues []Order

func (o Order) String() string {
	return fmt.Sprintf("Order[op=%s,id=%d "+
		"T[id=%d,p=%d, n=%d],"+
		"M[id=%d,p=%d, n=%d]", o.Op, o.Id, o.T.Id, o.T.Price, o.T.Num, o.M.Id, o.M.Price, o.M.Num)
}



//TODO 币币交易限制 （成交价-卖价）/卖价 < 30%
func MathQ(t, m int) bool {
	q := (float64(t) - float64(m)) / float64(m)
	r, _ := strconv.ParseFloat(strconv.FormatFloat(q, 'f', 2, 64), 64)
	//fmt.Println("mathq:", r)
	return math.Min(r, 0.3) == r
}

func (o *Order) PushOrder(op string, t Taker, m Maker) {
	//币币交易规则限制
	o.Lock.Lock()
	defer o.Lock.Unlock()
	push := Order{
		Id: orderId,
		Op: op,
		T:  t,
		M:  m,
	}
	//fmt.Println(push)
	OrderQueues = append(OrderQueues, push)
	orderId++
}
