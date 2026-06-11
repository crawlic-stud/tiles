package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeMove     MessageType = "move"
	MessageTypeRerender MessageType = "rerender"
)

type gameConn struct {
	gameID int64
	conn   *websocket.Conn
}

type gameMessage struct {
	gameID      int64
	messageType MessageType
	data        []byte
}

type GameHub struct {
	rooms      map[int64]map[*websocket.Conn]bool
	register   chan *gameConn
	unregister chan *gameConn
	broadcast  chan *gameMessage
}

func NewHub() *GameHub {
	return &GameHub{
		rooms:      make(map[int64]map[*websocket.Conn]bool),
		register:   make(chan *gameConn),
		unregister: make(chan *gameConn),
		broadcast:  make(chan *gameMessage),
	}
}

func (h *GameHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.closeAllRooms()
			log.Println("game hub shut down")
			return

		case conn := <-h.register:
			if h.rooms[conn.gameID] == nil {
				h.rooms[conn.gameID] = make(map[*websocket.Conn]bool)
			}
			h.rooms[conn.gameID][conn.conn] = true
			log.Printf("client connected to game %d, total conns: %d", conn.gameID, len(h.rooms[conn.gameID]))

		case conn := <-h.unregister:
			delete(h.rooms[conn.gameID], conn.conn)
			conn.conn.Close()
			log.Printf("client disconnected from game %d", conn.gameID)

		case message := <-h.broadcast:
			for conn := range h.rooms[message.gameID] {
				err := conn.WriteMessage(websocket.TextMessage, message.data)
				if err != nil {
					log.Println("write error:", err)
					conn.Close()
					delete(h.rooms[message.gameID], conn)
				}
			}
		}
	}
}

func (h *GameHub) closeAllRooms() {
	for gameID := range h.rooms {
		for conn := range h.rooms[gameID] {
			conn.Close()
		}
	}
}

func (h *GameHub) BroadcastJSON(gameID int64, messageType MessageType, message any) error {
	if h.rooms[gameID] == nil {
		return fmt.Errorf("no rooms found for game %d", gameID)
	}
	data, err := json.Marshal(map[string]any{
		"type": messageType,
		"data": message,
	})
	if err != nil {
		return err
	}
	h.broadcast <- &gameMessage{gameID, messageType, data}
	return nil
}

func (h *GameHub) BroadcastRerender(gameID int64) error {
	return h.BroadcastJSON(gameID, MessageTypeRerender, nil)
}

func (h *GameHub) Register(gameID int64, conn *websocket.Conn) {
	h.register <- &gameConn{gameID, conn}
}

func (h *GameHub) Unregister(gameID int64, conn *websocket.Conn) {
	h.unregister <- &gameConn{gameID, conn}
}
