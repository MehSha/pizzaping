package main

var ActiveOrders []*Order

func init() {
	ActiveOrders = make([]*Order, 0)
}
