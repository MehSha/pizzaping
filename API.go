package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Order struct {
	PizzaName string  `json:"pizza_name"`
	UserID    string  `json:"user_id"`
	Quantity  int     `json:"qty"`
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
}

func testOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
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
	restaurants, err := searchRestaurant(inp.PizzaName, inp.Lat, inp.Lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	fmt.Println("resaurant are:", restaurants)
	dishes := []*Dish{}
	for _, rest := range restaurants {
		d, err := getDish(rest.ID, inp.PizzaName)
		if err == nil {
			d.Restaurant = rest
			dishes = append(dishes, d)
		}
	}

	out, err := json.Marshal(dishes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func finalizeOrder(w http.ResponseWriter, r *http.Request) {

}
