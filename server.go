package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

var socketHub *Hub

func main() {
	fmt.Println("hi", os.Getenv("PORT"))

	db := connectDB()
	defer db.Close()

	initOrder(db)
	initRestaurant(db)

	socketHub = newHub()
	go socketHub.run()

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(socketHub, w, r)
	})
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	r.HandleFunc("/order", addOrderHr).Methods("POST")
	r.HandleFunc("/orders/{orderid}/accept", acceptOrderHr).Methods("POST")
	r.HandleFunc("/restaurant", registerRestaurantHr).Methods("POST")
	r.HandleFunc("/wstest", factorySendMessage(socketHub)).Methods("GET")
	r.HandleFunc("/wsusers", factoryGetUsers(socketHub)).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	srv := &http.Server{
		Handler: handlers.CORS(headersOk, originsOk, methodsOk)(r),
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("user")
	fmt.Println("user login, username:", userName)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if userName == "" {
				return false
			}
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), ID: userName}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
