package entity

import "github.com/google/uuid"

const MaxUsersPerRoom = 2

type Room struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	UsersId  []string   `json:"users"`
	Messages []*Message `json:"messages,omitempty"`
}

func NewRoom() *Room {
	return &Room{
		ID:      uuid.New().String(),
		Name:    uuid.New().String(),
		UsersId: make([]string, 0, MaxUsersPerRoom),
	}
}

func (t *Room) IsFull() bool {
	return len(t.UsersId) >= MaxUsersPerRoom
}

func (t *Room) IsEmpty() bool {
	return len(t.UsersId) == 0
}

func (t *Room) HasUser(userID string) bool {
	for _, id := range t.UsersId {
		if id == userID {
			return true
		}
	}
	return false
}

func (t *Room) AddUser(userId string) {
	if !t.HasUser(userId) {
		t.UsersId = append(t.UsersId, userId)
	}
}

func (t *Room) RemoveUser(userId string) {
	list := make([]string, 0, len(t.UsersId))
	for _, u := range t.UsersId {
		if u != userId {
			list = append(list, u)
		}
	}
	t.UsersId = list
}

func (t *Room) AddMessage(msg *Message) {
	t.Messages = append(t.Messages, msg)
}

func (t *Room) RemoveMessage(message *Message) {
	var list []*Message

	for _, msg := range t.Messages {
		if !msg.Equals(message.ID) {

			list = append(list, msg)
		}
	}

	t.Messages = list
}
