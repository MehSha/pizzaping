package main

import (
	"fmt"
	"log"

	"github.com/MehSha/basicdam"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type Order struct {
	ID        string  `json:"id"`
	PizzaName string  `json:"pizza_name"`
	UserID    string  `json:"user_id"`
	Status    string  `json:"status"`
	Address   string  `json:"address"`
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	RestaurantName string `json:"restaurant_name"`
}

var orderDAM *basicdam.BasicDAM

func initOrder(db *sqlx.DB) {
	orderDAM = basicdam.NewDAM(&Order{}, db)
	err := orderDAM.SyncDB()
	if err != nil {
		log.Fatalln(err)
	}
}

func AddOrder(ord *Order) (string,error) {
	ord.ID =  uuid.NewV4().String()
	ord.Status = "CREATED"
	id, err := orderDAM.Insert(ord)
	if err != nil {
		return "",err
	}
	log.Printf("order with ID: %s added", id)
	return ord.ID, nil
}

func acceptOrderDB(id, restName string) (*Order, error) {
	ord:= &Order{}
	err:=orderDAM.DB.Get(ord, "select * from orders where id=$1", id)
	if err != nil{
		log.Printf("could not get order:%s", err)
		return nil, err
	}
	if ord.Status != "CREATED" {
		return nil, fmt.Errorf("order already proccessed")
	}
	//update order
	_,err=orderDAM.DB.Exec("update orders set status='ACCEPTED', restaurantname=$2 where id=$1",  id, restName)
	ord.RestaurantName = restName
	ord.Status = "ACCEPTED"
	return ord, err
}
