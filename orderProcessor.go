package main

import "fmt"

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
		return fmt.Errorf("user %s does not have an active socket!\n", order.UserID)
	}
	socketHub.SendTo(order.UserID, []byte("order accepted"))
	return nil
}
