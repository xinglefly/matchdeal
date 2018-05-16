package model

import "fmt"

type Deal struct {
	Price int
	Num   int
}

func Create(price, num int) *Deal {
	return &Deal{
		Price: price,
		Num:   num,
	}
}

func (d *Deal) SetNum(num int){
	if d == nil {
		fmt.Println("Deal is nil isIgnore.")
		return
	}
	d.Num = num
}
