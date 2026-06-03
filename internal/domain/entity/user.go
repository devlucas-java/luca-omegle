package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID        string
	UserName  string
	IsWaiting bool
	IsOnline  bool
	RoomID    string
}

func NewUser(username string) *User {
	return &User{
		ID:       uuid.New().String(),
		UserName: username,
	}
}

func (t *User) Equals(id string) bool {
	return t.ID == id
}

func (u *User) AssignRoom(roomID string) {
	u.RoomID = roomID
	u.IsWaiting = false
}

func (u *User) ClearRoom() {
	u.RoomID = ""
	u.IsWaiting = true
}
