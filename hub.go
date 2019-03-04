// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients   map[*Client]bool
	clientMap map[string]*Client
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		clientMap:  make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientMap[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientMap, client.ID)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) SendTo(id string, message []byte) {
	cl, ok := h.clientMap[id]
	if !ok {
		return
	}
	cl.send <- message
}

func (h *Hub) SendGroup(ids []string, message []byte) {
	for _, id := range ids {
		h.SendTo(id, message)
	}
}

func factorySendMessage(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := strings.Split(r.URL.Query().Get("id"), ",")
		message := []byte(r.URL.Query().Get("message"))
		hub.SendGroup(ids, message)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok.\n")
	}
}
func factoryGetUsers(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := strings.Split(r.URL.Query().Get("id"), ",")
		message := []byte(r.URL.Query().Get("message"))
		hub.SendGroup(ids, message)
		users := []string{}
		for u := range hub.clientMap {
			users = append(users, u)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v\n", users)
	}
}
