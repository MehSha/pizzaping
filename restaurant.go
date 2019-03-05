package main

import (
	"log"

	"github.com/MehSha/basicdam"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	webpush "github.com/SherClockHolmes/webpush-go"
)

type RestaurantInput struct {
	ID           string `json:"id"`
	Name         string    `json:"name"`
	Subscription webpush.Subscription `json:"subscription"`
}

type Restaurant struct {
	ID           string `json:"id"`
	Name         string    `json:"name"`
	EndPoint string `json:"endpoint"`
	P256DH string `json:"p256dh"`
	Auth string `json:"auth"`
}

var restaurantDAM *basicdam.BasicDAM

func initRestaurant(db *sqlx.DB) {
	restaurantDAM = basicdam.NewDAM(&Restaurant{}, db)
	err := restaurantDAM.SyncDB()
	if err != nil {
		log.Fatalln(err)
	}
}

func AddRestaurant(rest *Restaurant) error {
	rest.ID =  uuid.NewV4().String()
	id, err := restaurantDAM.Insert(rest)
	if err != nil {
		return err
	}
	log.Printf("restaurant with ID: %s added", id)
	return nil
}

func ListRestaurants()([]Restaurant, error){
	result:= []Restaurant{}
	err:= restaurantDAM.DB.Select(&result, "select * from restaurants")
	return result, err
}
