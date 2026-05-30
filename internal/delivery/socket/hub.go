package socket

import (
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	model "github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type Hub struct {
	Room map[string]*dto.Room
	log  *logger.Logger
}

func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		Room: make(map[string]*dto.Room),
		log:  log,
	}
}

func (t *Hub) RegisterRoom(c *model.Client) {
}

func (t *Hub) UnregisterRoom(c *model.Client) {

}
