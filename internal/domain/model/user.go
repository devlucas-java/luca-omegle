package model

import (
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Conn     *websocket.Conn
}

func NewUser(username string, conn *websocket.Conn) *User {
	return &User{
		ID:       uuid.New().String(),
		UserName: username,
		Conn:     conn,
	}
}

func (t *User) Equals(id string) bool {
	return t.ID == id
}
