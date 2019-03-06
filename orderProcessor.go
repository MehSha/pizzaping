package main

import (
	"encoding/json"
	"log"
)

var ActiveOrders []*Order

func init() {
	ActiveOrders = make([]*Order, 0)
}

func acceptOrder(id, restName string) error {
	//update DB
	order, err := acceptOrderDB(id, restName)
	if err != nil {
		return err
	}
	//send message to user
	_, ok := socketHub.clientMap[order.UserID]
	if !ok {
		log.Printf("user %s does not have an active socket!\n", order.UserID)
	} else {
		bt, _ := json.Marshal(order)
		socketHub.SendTo(order.UserID, bt)
	}

	return nil
}
