package socket

import (
	"context"
	"fmt"
	"sync"

	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type sendEvent struct {
	session *dto.Session
	frame   *dto.Frame
}

type subscribeEvent struct {
	session *dto.Session
}

type unsubscribeEvent struct {
	session *dto.Session
}

type disconnectEvent struct {
	session *dto.Session
}

// Hub manages all active sessions, the waiting queue, and room pairing.
type Hub struct {
	log         *logger.Logger
	roomService *service.RoomService
	userService *service.UserService

	sessions sync.Map // userID → *dto.Session

	sendCh        chan sendEvent
	subscribeCh   chan subscribeEvent
	unsubscribeCh chan unsubscribeEvent
	disconnectCh  chan disconnectEvent
	waitingQueue  chan *dto.Session
}

func NewHub(
	log *logger.Logger,
	roomService *service.RoomService,
	userService *service.UserService,
) *Hub {
	return &Hub{
		log:           log,
		roomService:   roomService,
		userService:   userService,
		sendCh:        make(chan sendEvent, 256),
		subscribeCh:   make(chan subscribeEvent, 256),
		unsubscribeCh: make(chan unsubscribeEvent, 256),
		disconnectCh:  make(chan disconnectEvent, 256),
		waitingQueue:  make(chan *dto.Session, 256),
	}
}

// Run starts the hub event loop and the matchmaker goroutine.
// It blocks until ctx is cancelled.
func (h *Hub) Run(ctx context.Context) {
	go h.matchmaker(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case ev := <-h.sendCh:
			go h.handleSend(ctx, ev)

		case ev := <-h.subscribeCh:
			go h.handleSubscribe(ctx, ev)

		case ev := <-h.unsubscribeCh:
			go h.handleUnsubscribe(ctx, ev)

		case ev := <-h.disconnectCh:
			go h.handleDisconnect(ctx, ev)
		}
	}
}

func (h *Hub) matchmaker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case s1 := <-h.waitingQueue:
			select {
			case <-ctx.Done():
				return
			case s2 := <-h.waitingQueue:
				go h.pair(ctx, s1, s2)
			}
		}
	}
}

func (h *Hub) pair(ctx context.Context, s1, s2 *dto.Session) {
	room, err := h.roomService.JoinRandom(ctx, []string{s1.User.ID, s2.User.ID})
	if err != nil {
		h.log.Errorf("pair: %v", err)
		_ = s1.SendError(ctx, "matchmaking failed")
		_ = s2.SendError(ctx, "matchmaking failed")
		return
	}

	// JoinRandom updates RoomID inside the user structs passed to it,
	// but s1.User / s2.User are the original pointers — sync them back.
	s1.User.AssignRoom(room.ID)
	s2.User.AssignRoom(room.ID)

	connFrame := dto.NewFrame(dto.CmdConnected, "", map[string]string{
		"room-id": room.ID,
	})
	if err := s1.Send(ctx, connFrame); err != nil {
		h.log.Errorf("pair send s1: %v", err)
	}
	if err := s2.Send(ctx, connFrame); err != nil {
		h.log.Errorf("pair send s2: %v", err)
	}

	h.log.Infof("paired users %s and %s in room %s", s1.User.ID, s2.User.ID, room.ID)
}

func (h *Hub) handleSend(ctx context.Context, ev sendEvent) {
	roomID, err := ev.frame.MustHeader("room-id")
	if err != nil {
		_ = ev.session.SendError(ctx, err.Error())
		return
	}

	out := dto.NewFrame(dto.CmdMessage, ev.frame.Body, map[string]string{
		"room-id":  roomID,
		"username": ev.session.User.UserName,
	})

	h.sessions.Range(func(_, val any) bool {
		s, ok := val.(*dto.Session)
		if !ok {
			return true
		}
		if s.User.RoomID == roomID && s.User.ID != ev.session.User.ID {
			if err := s.Send(ctx, out); err != nil {
				h.log.Errorf("handleSend to %s: %v", s.User.ID, err)
			}
		}
		return true
	})
}

func (h *Hub) handleSubscribe(ctx context.Context, ev subscribeEvent) {
	if ev.session.User.RoomID != "" {
		if err := h.roomService.LeaveRoom(ctx, ev.session.User); err != nil {
			h.log.Errorf("handleSubscribe leaveRoom: %v", err)
			_ = ev.session.SendError(ctx, fmt.Sprintf("could not leave room: %v", err))
			return
		}
		ev.session.User.ClearRoom()
	}
	h.waitingQueue <- ev.session
	h.log.Infof("user %s added to waiting queue", ev.session.User.ID)
}

func (h *Hub) handleUnsubscribe(ctx context.Context, ev unsubscribeEvent) {
	if ev.session.User.RoomID == "" {
		return
	}
	if err := h.roomService.LeaveRoom(ctx, ev.session.User); err != nil {
		h.log.Errorf("handleUnsubscribe: %v", err)
		_ = ev.session.SendError(ctx, fmt.Sprintf("could not leave room: %v", err))
		return
	}
	ev.session.User.ClearRoom()
}

func (h *Hub) handleDisconnect(ctx context.Context, ev disconnectEvent) {
	if ev.session.User.RoomID != "" {
		if err := h.roomService.LeaveRoom(ctx, ev.session.User); err != nil {
			h.log.Errorf("handleDisconnect leaveRoom: %v", err)
		}
		ev.session.User.ClearRoom()
	}
	if err := h.userService.Delete(ctx, ev.session.User); err != nil {
		h.log.Errorf("handleDisconnect deleteUser: %v", err)
	}
	h.sessions.Delete(ev.session.User.ID)
	h.log.Infof("user %s disconnected", ev.session.User.ID)
}

func (h *Hub) Register(s *dto.Session) {
	h.sessions.Store(s.User.ID, s)
}

func (h *Hub) Unregister(s *dto.Session) {
	h.sessions.Delete(s.User.ID)
}
