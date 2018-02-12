package sharing

import (
	"github.com/gavincabbage/api.il2missionplanner.com/models"
	"go/doc"
)

type Room struct {
	Clients     map[*Client]bool
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	MissionPlan *models.Plan
}

func NewRoom() *Room {
	return &Room{
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
		MissionPlan: &models.Plan{},
	}
}

func (h *Room) Start() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
