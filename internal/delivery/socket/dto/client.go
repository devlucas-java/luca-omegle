package dto

import (
	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
)

type Session struct {
	Conn *websocket.Conn
	User *entity.User
}

func NewSession(c *websocket.Conn, u *entity.User) *Session {
	return &Session{
		Conn: c,
		User: u,
	}
}

type Broadcast struct {
	*Session
	Msg *entity.Message
}

func NewBroadcast(s *Session, m *entity.Message) *Broadcast {
	return &Broadcast{
		Session: s,
		Msg:     m,
	}
}
