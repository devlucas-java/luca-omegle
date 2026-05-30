package main

import (
	"net/http"

	_ "github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/configs"
	"github.com/devlucas-java/luca-omegle/pkg/logger"

	"github.com/devlucas-java/luca-omegle/internal/delivery/socket"
	_ "github.com/panjf2000/gnet/v2"
	_ "github.com/pion/webrtc/v4"
)

func main() {

	conf := configs.InitConfig()
	log := logger.NewLogger(logger.TRACE)

	s := socket.NewWSHandler(log.WithComponent("WS_Handler"))
	h := socket.NewDashboard(log.WithComponent("WS_Dashboard"))

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.WSHandlerS)
	mux.HandleFunc("/get", h.GetClients)
	mux.HandleFunc("/ws/dashboard", h.WSDashboard)
	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	port := ":" + conf.ServerPort
	http.ListenAndServe(port, mux)
}
