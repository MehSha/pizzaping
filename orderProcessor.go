package main

import "log"

var ActiveOrders []*Order

func init() {
	ActiveOrders = make([]*Order, 0)
}

func acceptOrder(id string) error {
	//update DB
	order, err := acceptOrderDB(id)
	if err != nil {
		return err
	}
	//send message to user
	_, ok := socketHub.clientMap[order.UserID]
	if !ok {
		log.Printf("user %s does not have an active socket!\n", order.UserID)
	} else {
		socketHub.SendTo(order.UserID, []byte("order accepted"))
	}

	return nil
}
