package dto

import (
	"github.com/coder/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	Message  chan *Message
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}
