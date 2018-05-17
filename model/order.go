package model

type Order struct {
	Id int
	Op string
	T  Taker
	M  Maker
}


var OrderQueues []Order