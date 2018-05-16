package main

import (
	"matchdeal/model"
)

//TODO 后期重构用接口
type Goruting interface {
	CreateGoruting() chan int
	Receiver(id int) chan<- int
}

func main() {
	model.CreateTaker()
	model.CreateMaker()
}
