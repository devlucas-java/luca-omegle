package socket

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type WSHandler struct {
	log             *logger.Logger
	userService     *service.UserService
	broadcastChan   chan dto.Session
	subscribeChan   chan dto.Session
	unSubscribeChan chan dto.Session
	dashboardChan   chan dto.Session
	disconnectChan  chan dto.Session
}

func NewWSHandler(
	l *logger.Logger,
	us *service.UserService,
) *WSHandler {
	return &WSHandler{
		log:             l,
		userService:     us,
		broadcastChan:   make(chan dto.Session, 100),
		subscribeChan:   make(chan dto.Session, 100),
		unSubscribeChan: make(chan dto.Session, 100),
		dashboardChan:   make(chan dto.Session, 100),
		disconnectChan:  make(chan dto.Session, 100),
	}
}

func (t *WSHandler) WSHandlerS(w http.ResponseWriter, r *http.Request) {

	nickname := getNickname(r)

	conn, err := acceptWS(w, r)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusAbnormalClosure, "closed")

	for {

		_, data, _ := conn.Read(r.Context())

		switch string(data) {
		case "BROADCAST":
			go Broadcast(nickname, string(data))

		case "SUBSCRIBE":
			go Subscribe(nickname)

		case "UNSUBSCRIBE":
			go Unsubscribe(nickname)

		case "DASHBOARD":
			go Dashboard()

		case "DISCONNECT":
			go Disconnect(nickname)
		}
	}
}

func acceptWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns:     []string{"*"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getNickname(r *http.Request) string {
	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		nickname = "anonymous"
	}
	return nickname
}
