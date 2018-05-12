// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

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
	}
}

func (hubinstance *Hub) run() {
	for {
		select {
		case client := <-hubinstance.register:
			hubinstance.clients[client] = true
		case client := <-hubinstance.unregister:
			if _, ok := hubinstance.clients[client]; ok {
				delete(hubinstance.clients, client)
				close(client.send)
			}
		case message := <-hubinstance.broadcast:
			for client := range hubinstance.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hubinstance.clients, client)
				}
			}
		}
	}
}
