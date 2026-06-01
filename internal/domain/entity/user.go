package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID        string
	UserName  string
	IsWaiting bool
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
