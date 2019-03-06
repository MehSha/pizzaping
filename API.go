package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func testOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

type OrderNotif struct {
	OrderId string `json:"order_id"`
	Address string `json:"address"`
	Dish    string `json:"dish"`
}

func addOrderHr(w http.ResponseWriter, r *http.Request) {
	//get request object
	inp := &Order{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, inp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// got order, now first persist it!
	orderid, err := AddOrder(inp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("coud not persist order", err)
		return
	}
	//get list of restaurants and push message to them
	rests, err := ListRestaurants()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("coud not list restaurants", err)
		return
	}
	fmt.Printf("restaurants are: %+v\n", rests)
	for _, rst := range rests {
		notif := OrderNotif{
			OrderId: orderid,
			Dish:    inp.PizzaName,
			Address: inp.Address,
		}
		notifMsg, _ := json.Marshal(notif)
		pushNotif(rst, string(notifMsg))
	}
	//done!
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func acceptOrderHr(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["orderid"]
	err := acceptOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("coud not update order", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func registerRestaurantHr(w http.ResponseWriter, r *http.Request) {
	//get request object
	inp := &RestaurantInput{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error reading request body", err)
		return
	}
	err = json.Unmarshal(body, inp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error unmarshal request body", err)
		return
	}
	if inp.Subscription.Endpoint == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "restaurant must have subscription information\n")
		return
	}
	// got restaurant, now first persist it!
	rst := &Restaurant{
		Name:     inp.Name,
		EndPoint: inp.Subscription.Endpoint,
		P256DH:   inp.Subscription.Keys.P256dh,
		Auth:     inp.Subscription.Keys.Auth,
	}
	err = AddRestaurant(rst)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// pushNotif(*rst, "you are added!")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")

}
func finalizeOrder(w http.ResponseWriter, r *http.Request) {

}
