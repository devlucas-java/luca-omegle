package dto

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
)

type Session struct {
	Conn *websocket.Conn
	User *entity.User
}

func NewSession(c *websocket.Conn, u *entity.User) *Session {
	return &Session{Conn: c, User: u}
}

func (s *Session) Send(ctx context.Context, f *Frame) error {
	data := f.Encode()
	if err := s.Conn.Write(ctx, websocket.MessageText, []byte(data)); err != nil {
		return fmt.Errorf("session %s write: %w", s.User.ID, err)
	}
	return nil
}

func (s *Session) SendError(ctx context.Context, msg string) error {
	f := NewFrame(CmdError, msg, map[string]string{"message": msg})
	return s.Send(ctx, f)
}
