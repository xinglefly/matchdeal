package model

import (
	"strconv"
	"math"
	"fmt"
)

type Order struct {
	Id int
	Op string
	T  Taker
	M  Maker
}

var OrderQueues []Order

//TODO 币币交易限制 （成交价-卖价）/卖价 < 30%
func MathQ(t, m int) bool {
	q := (float64(t) - float64(m)) / float64(m)
	r, _ := strconv.ParseFloat(strconv.FormatFloat(q, 'f', 2, 64), 64)
	fmt.Println("mathq:", r)
	return math.Min(r, 0.3) == r
}
