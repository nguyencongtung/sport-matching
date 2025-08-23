package handler

import (
	"log"
	"sync"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// RoomManager manages chat rooms and their clients
type RoomManager struct {
	clients map[string]map[*websocket.Conn]bool // Room -> Clients
	mutex   sync.RWMutex                        // Mutex for thread-safe access
}

var roomManager = &RoomManager{
	clients: make(map[string]map[*websocket.Conn]bool),
}

// ChatHandler handles WebSocket connections for chat
func ChatHandler(c *fiber.Ctx) error {
	// Only upgrade to WebSocket if the request is a WebSocket request
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// WebSocketHandler handles the WebSocket connection
func WebSocketHandler(c *websocket.Conn) {
	room := c.Query("room")
	if room == "" {
		slog.Info("Missing room parameter")
		c.Close()
		return
	}

	roomManager.registerClient(room, c)
	defer func() {
		roomManager.unregisterClient(room, c)
		c.Close()
	}()

	slog.Info("New client connected to room: %s\n", room)

	for {
		var msg map[string]string
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		slog.Info("Message received in room %s: %v\n", room, msg)
		roomManager.broadcastMessage(room, msg, c)
	}
}

// registerClient adds a client to a room
func (rm *RoomManager) registerClient(room string, conn *websocket.Conn) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if rm.clients[room] == nil {
		rm.clients[room] = make(map[*websocket.Conn]bool)
	}
	rm.clients[room][conn] = true
	slog.Info("Client registered in room: %s\n", room)
}

// unregisterClient removes a client from a room
func (rm *RoomManager) unregisterClient(room string, conn *websocket.Conn) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if rm.clients[room] != nil {
		delete(rm.clients[room], conn)
		if len(rm.clients[room]) == 0 {
			delete(rm.clients, room)
		}
	}
	slog.Info("Client unregistered from room: %s\n", room)
}

// broadcastMessage sends a message to all clients in a room except the sender
func (rm *RoomManager) broadcastMessage(room string, msg map[string]string, sender *websocket.Conn) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	for client := range rm.clients[room] {
		if client != sender {
			err := client.WriteJSON(msg)
			if err != nil {
				slog.Info("Error broadcasting message:", err)
				client.Close()
				rm.unregisterClient(room, client)
			}
		}
	}
}
