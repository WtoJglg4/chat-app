package models

import "log"

type Hub struct {
	clients    map[*Client]struct{}
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Register(client *Client) {
	hub.register <- client
}

func (hub *Hub) Unregister(client *Client) {
	hub.unregister <- client
}

func (hub *Hub) SendAll(msg []byte) {
	hub.broadcast <- msg
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = struct{}{}
			log.Printf("hub: register: %v\n", client.Username)
		case client := <-hub.unregister:
			delete(hub.clients, client)
			close(client.Send)
			log.Printf("hub: unregister: %v\n", client.Username)
		case msg := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.Send <- msg:
				default:
					log.Printf("hub: client %v is inactive\n", client.Username)
					delete(hub.clients, client)
					close(client.Send)
				}
			}
			log.Printf("hub: broadcast message: %v\n", string(msg))
		}
	}
}

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}
