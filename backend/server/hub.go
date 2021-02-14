package server

import (
	"EdgeNews/backend/utils"
	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and
// broadcasts messages to the clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

var newsHub *Hub
var textStreamHub *Hub

func init() {
	log.Info("[HUB] Creating hubs...")
	newsHub = NewHub()
	textStreamHub = NewHub()

	go newsHub.RunForever()
	go textStreamHub.RunForever()
}

func GetNewsHub() *Hub {
	return newsHub
}

func GetTextStreamHub() *Hub {
	return textStreamHub
}

func (h *Hub) SendMessage(message []byte) {
	log.Info("[HUB] Broadcast: " + string(message))
	h.broadcast <- message
}

func (h *Hub) RunForever() {
	log.Debug("[HUB] Starting...")

	for {
		select {
		case client := <-h.register:
			log.Debug("[HUB] Register new client")
			utils.GetMetrics().WS.ConnectionsActive.Inc()
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Debug("[HUB] Unregister new client")
				utils.GetMetrics().WS.ConnectionsActive.Dec()
				delete(h.clients, client)
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
