package main

import (
	"net/http"
	"time"

	_ "github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/configs"
	"github.com/devlucas-java/luca-omegle/pkg/logger"

	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/cache"
	_ "github.com/panjf2000/gnet/v2"
	_ "github.com/pion/webrtc/v4"
)

func main() {

	conf := configs.InitConfig()
	log := logger.NewLogger(logger.TRACE)
	ch := configs.InitCache(conf)

	userRepository := cache.NewUserCH(ch, 20*time.Minute)
	userService := service.NewUserService(userRepository, log.WithComponent("User_Service"))
	s := socket.NewWSHandler(log.WithComponent("WS_Handler"), userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.WSHandlerS)
	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	port := ":" + conf.ServerPort
	http.ListenAndServe(port, mux)
}
