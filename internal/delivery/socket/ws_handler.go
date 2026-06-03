package socket

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"
	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type WSHandler struct {
	log           *logger.Logger
	userService   *service.UserService
	broadcastChan chan *dto.Broadcast
	waitingChan   chan *dto.Session
}

func NewWSHandler(
	l *logger.Logger,
	us *service.UserService,
) *WSHandler {
	return &WSHandler{
		log:           l,
		userService:   us,
		broadcastChan: make(chan *dto.Broadcast, 100),
		waitingChan:   make(chan *dto.Session, 100),
	}
}

func (t *WSHandler) WSHandlerS(w http.ResponseWriter, r *http.Request) {

	nickname := getNickname(r)

	conn, err := acceptWS(w, r)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusAbnormalClosure, "closed")

	user, err := t.userService.Register(r.Context(), entity.NewUser(nickname))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session := dto.NewSession(conn, user)
	for {

		_, data, _ := conn.Read(r.Context())

		switch string(data) {
		case "BROADCAST":
			go Broadcast(nickname, string(data)) // tengo que hacer una funcion de extrair el type de la message y extrair el contenido
			msg := entity.NewMessage("", "", nickname)
			t.broadcastChan <- dto.NewBroadcast(session, msg)

		case "SUBSCRIBE":
			go Subscribe(nickname)

		case "UNSUBSCRIBE":
			go Unsubscribe(nickname)
			t.waitingChan <- session

		case "DASHBOARD":
			go Dashboard()

		case "DISCONNECT":
			Disconnect(nickname, session)
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
