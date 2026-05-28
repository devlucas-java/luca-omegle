package socket

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type Server struct {
	log *logger.Logger
}

func NewServer(l *logger.Logger) *Server {
	return &Server{log: l}
}

var (
	clients map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
)

func (t *Server) HandlerWS(w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns:     []string{"*"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.log.Error("error in accept web socket: ", err.Error())
		return
	}

	clients[conn] = true
	t.log.Info("new connection")

	defer conn.Close(websocket.StatusAbnormalClosure, "closed")

	for {

		msgType, data, err := conn.Read(r.Context())
		if err != nil {
			clients[conn] = false
			t.log.Error("error in read web socket: ", err.Error())
			break
		}

		t.log.Info("message type: ", msgType)
		t.log.Info("data: ", string(data))

		for client, b := range clients {
			if b {
				err = client.Write(r.Context(), websocket.MessageText, data)
				if err != nil {
					clients[conn] = false
					t.log.Error("error in write web socket: ", err.Error())
					break
				}
			}
		}
	}
}

func (t *Server) GetClients(w http.ResponseWriter, r *http.Request) {

	var clientsTrue int

	for _, v := range clients {
		if v {
			clientsTrue++
			fmt.Println(clientsTrue)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(clientsTrue)))
}
