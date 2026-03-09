package websocket

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Map of user ID to their active clients (allows multiple tabs/devices)
	userClients map[int]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// GetOnlineUserIDs returns a slice of currently online user IDs
func (h *Hub) GetOnlineUserIDs() []int {
	ids := make([]int, 0, len(h.userClients))
	for uid := range h.userClients {
		ids = append(ids, uid)
	}
	return ids
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[int]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.userClients[client.userID]; !ok {
				h.userClients[client.userID] = make(map[*Client]bool)
			}
			h.userClients[client.userID][client] = true
			log.Printf("WebSocket: user %d connected (client added)", client.userID)
			h.broadcastUserStatus(client.userID, true)

			// Send current presence list to the newly registered client
			userIDs := make([]int, 0, len(h.userClients))
			for uid := range h.userClients {
				userIDs = append(userIDs, uid)
			}
			presenceMsg := map[string]interface{}{
				"type":     "presence_list",
				"user_ids": userIDs,
			}
			jsonPresence, err := json.Marshal(presenceMsg)
			if err != nil {
				log.Println("Error marshaling presence list:", err)
			} else {
				select {
				case client.send <- jsonPresence:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				if userClients, ok := h.userClients[client.userID]; ok {
					delete(userClients, client)
					if len(userClients) == 0 {
						delete(h.userClients, client.userID)
						log.Printf("WebSocket: user %d fully disconnected (no client left)", client.userID)
						h.broadcastUserStatus(client.userID, false)
					}
				}
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

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID int, message []byte) {
	if clients, ok := h.userClients[userID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(clients, client)
			}
		}
	}
}

func (h *Hub) broadcastUserStatus(userID int, isOnline bool) {
	status := "offline"
	if isOnline {
		status = "online"
	}

	msg := map[string]interface{}{
		"type":    "user_status",
		"user_id": userID,
		"status":  status,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling status message:", err)
		return
	}

	log.Printf("WebSocket: broadcasting user_status - userID=%d status=%s to %d clients", userID, status, len(h.clients))

	// Broadcast to everyone
	// In a real app, might only broadcast to friends/connections
	for client := range h.clients {

		select {
		case client.send <- jsonMsg:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
