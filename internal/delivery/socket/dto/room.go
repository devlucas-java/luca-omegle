package dto

import "github.com/google/uuid"

type Room struct {
	ID      uuid.UUID          `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func NewRoom(name string) *Room {
	return &Room{
		ID:      uuid.New(),
		Name:    name,
		Clients: make(map[string]*Client),
	}
}

func (t *Room) AddClient(c *Client) {
	t.Clients[c.ID] = c
}

func (t *Room) RemoveClient(c *Client) {
	delete(t.Clients, c.ID)
}

func (t *Room) GetClient(id string) *Client {
	return t.Clients[id]
}
