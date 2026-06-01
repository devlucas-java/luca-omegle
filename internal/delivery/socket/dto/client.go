package dto

import (
	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
)

type Session struct {
	Conn *websocket.Conn
	User *entity.User
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}
