package main

import (
	"net/http"

	_ "github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/pkg/logger"

	"github.com/devlucas-java/luca-omegle/internal/delivery/socket"
	_ "github.com/panjf2000/gnet/v2"
	_ "github.com/pion/webrtc/v4"
)

func main() {

	s := socket.NewServer(logger.NewLogger(logger.TRACE))
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandlerWS)
	mux.HandleFunc("/get", s.GetClients)
	http.ListenAndServe(":8080", mux)
}
