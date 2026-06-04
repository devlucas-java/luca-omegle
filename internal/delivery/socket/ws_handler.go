package socket

import (
	"errors"
	"net/http"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type WSHandler struct {
	log         *logger.Logger
	userService *service.UserService
	hub         *Hub
}

func NewWSHandler(log *logger.Logger, us *service.UserService, hub *Hub) *WSHandler {
	return &WSHandler{log: log, userService: us, hub: hub}
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := acceptWS(w, r)
	if err != nil {
		h.log.Errorf("ws accept: %v", err)
		return
	}

	nickname := nickFromQuery(r)
	user, err := h.userService.Register(r.Context(), entity.NewUser(nickname))
	if err != nil {
		h.log.Errorf("register user: %v", err)
		conn.Close(websocket.StatusInternalError, "registration failed")
		return
	}

	session := dto.NewSession(conn, user)
	h.hub.Register(session)
	defer func() {
		h.hub.disconnectCh <- disconnectEvent{session: session}
		conn.Close(websocket.StatusNormalClosure, "bye")
	}()

	h.log.Infof("user %s (%s) connected", user.ID, nickname)

	for {
		_, data, err := conn.Read(r.Context())
		if err != nil {
			var wsErr websocket.CloseError
			if errors.As(err, &wsErr) {
				h.log.Infof("user %s closed connection: %v", user.ID, wsErr)
			} else {
				h.log.Errorf("user %s read error: %v", user.ID, err)
			}
			return
		}

		frame, err := dto.ParseFrame(string(data))
		if err != nil {
			h.log.Warnf("user %s bad frame: %v", user.ID, err)
			_ = session.SendError(r.Context(), "bad frame: "+err.Error())
			continue
		}

		h.dispatch(session, frame)
	}
}

func (h *WSHandler) dispatch(s *dto.Session, f *dto.Frame) {
	switch f.Command {
	case dto.CmdSend:
		h.hub.sendCh <- sendEvent{session: s, frame: f}

	case dto.CmdSubscribe:
		h.hub.subscribeCh <- subscribeEvent{session: s}

	case dto.CmdNext:
		h.hub.subscribeCh <- subscribeEvent{session: s}

	case dto.CmdDisconnect:
		h.hub.disconnectCh <- disconnectEvent{session: s}

	case dto.CmdDashboard:
		h.log.Debugf("DASHBOARD requested by %s", s.User.ID)

	default:
		h.log.Warnf("unknown command %q from %s", f.Command, s.User.ID)
		_ = s.SendError(nil, "unknown command: "+string(f.Command))
	}
}

func acceptWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns:     []string{"*"},
		InsecureSkipVerify: true,
	})
}

func nickFromQuery(r *http.Request) string {
	if n := r.URL.Query().Get("nickname"); n != "" {
		return n
	}
	return "anonymous"
}
