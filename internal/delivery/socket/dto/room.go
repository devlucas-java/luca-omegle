package dto

import "github.com/google/uuid"

type Room struct {
	ID      uuid.UUID           `json:"id"`
	Name    string              `json:"name"`
	Clients map[string]*Session `json:"clients"`
}

func NewRoom(name string) *Room {
	return &Room{
		ID:      uuid.New(),
		Name:    name,
		Clients: make(map[string]*Session),
	}
}

func (t *Room) AddClient(c *Session) {
	t.Clients[c.User.ID] = c
}

func (t *Room) RemoveClient(c *Session) {
	delete(t.Clients, c.User.ID)
}

func (t *Room) GetClient(id string) *Session {
	return t.Clients[id]
}
