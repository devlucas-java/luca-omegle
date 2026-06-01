package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	RoomID    string `json:"room_id"`
	Username  string `json:"username"`
}

func NewMessage(content string, roomID string, username string) *Message {
	return &Message{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
		RoomID:    roomID,
		Username:  username,
	}
}

func (t *Message) Equals(id string) bool {
	return t.ID == id
}
