package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hi", os.Getenv("PORT"))

	db := connectDB()
	defer db.Close()

	initProduct(db)

	hub := newHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	r.HandleFunc("/order", orderHandler).Methods("POST")
	r.HandleFunc("/wstest", factorySendMessage(hub)).Methods("GET")
	r.HandleFunc("/wsusers", factoryGetUsers(hub)).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	fmt.Println("user login, email:", email)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if email == "" {
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
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), ID: email}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
