package socket

import (
	"net/http"
	"strconv"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type WSHandler struct {
	log *logger.Logger
}

func NewWSHandler(l *logger.Logger) *WSHandler {
	return &WSHandler{log: l}
}

type Client struct {
	Nickname string          `json:"nickname"`
	conn     *websocket.Conn `json:"-"`
}

var (
	clients map[*Client]bool = make(map[*Client]bool)
)

func (t *WSHandler) WSHandlerS(w http.ResponseWriter, r *http.Request) {

	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		nickname = "anonymous"
	}
	t.log.Info("nickname: ", nickname)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns:     []string{"*"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.log.Error("error in accept web socket: ", err.Error())
		return
	}

	client := &Client{
		Nickname: nickname,
		conn:     conn,
	}

	clients[client] = true
	t.log.Info("new connection" + nickname + " " + strconv.Itoa(len(clients)))

	defer conn.Close(websocket.StatusAbnormalClosure, "closed")

	for c, b := range clients {
		if b {
			err = c.conn.Write(r.Context(), websocket.MessageText, []byte(nickname+" connected"))
			if err != nil {
				clients[client] = false
				t.log.Error("error in write web socket: ", err.Error())
				break
			}
		}
	}

	for {

		msgType, data, err := conn.Read(r.Context())
		if err != nil {
			clients[client] = false
			t.log.Error("error in read web socket: ", err.Error())
			break
		}

		t.log.Info("message type: ", msgType)
		t.log.Info("data: ", string(data))

		for c, b := range clients {
			if b {
				err = c.conn.Write(r.Context(), websocket.MessageText, []byte(nickname+": "+string(data)))
				if err != nil {
					clients[client] = false
					t.log.Error("error in write web socket: ", err.Error())
					break
				}
			}
		}
	}
}
