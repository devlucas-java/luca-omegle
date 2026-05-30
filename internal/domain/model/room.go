package model

import "github.com/google/uuid"

type Room struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Users    []*User `json:"users"`
	Messages []*Message
}

func NewRoom(name string) *Room {
	return &Room{
		ID:    uuid.New().String(),
		Name:  name,
		Users: make([]*User, 0),
	}
}

func (t *Room) AddUser(user *User) {
	t.Users = append(t.Users, user)
}

func (t *Room) RemoveUser(user *User) {
	var list []*User

	for _, u := range t.Users {
		if !u.Equals(user.ID) {

			list = append(list, u)
			t.Users = list
		}
	}
}

func (t *Room) AddMessage(msg *Message) {
	t.Messages = append(t.Messages, msg)
}

func (t *Room) RemoveMessage(message *Message) {
	var list []*Message

	for _, msg := range t.Messages {
		if !msg.Equals(message.ID) {

			list = append(list, msg)
			t.Messages = list
		}
	}
}
