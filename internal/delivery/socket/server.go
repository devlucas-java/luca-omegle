package socket

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type Server struct {
	log *logger.Logger
}

func NewServer(l *logger.Logger) *Server {
	return &Server{log: l}
}

func (t *Server) HandlerWS(w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		t.log.Error("error in accept web socket: ", err.Error())
		return
	}
	t.log.Info("new connection")

	defer conn.Close(websocket.StatusAbnormalClosure, "closed")

	ctx := context.Background()

	for {

		msgType, data, err := conn.Read(ctx)
		if err != nil {
			t.log.Error("error in read web socket: ", err.Error())
		}

		t.log.Info("message type: ", msgType)
		t.log.Info("data: ", string(data))

		err = conn.Write(ctx, msgType, data)
		if err != nil {
			t.log.Error("error in write web socket: ", err.Error())
		}
	}
}
