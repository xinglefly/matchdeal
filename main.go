package main

import (
	"matchdeal/model"
)

type Goruting interface {
	CreateGoruting() chan int
	Receiver(id int) chan<- int
}

func main() {
	model.CreateTaker()
	model.CreateMaker()
}
