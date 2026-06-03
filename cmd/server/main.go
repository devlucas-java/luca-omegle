package main

import (
	"net/http"
	"time"

	"github.com/devlucas-java/luca-omegle/configs"
	"github.com/devlucas-java/luca-omegle/internal/application/service"
	"github.com/devlucas-java/luca-omegle/internal/delivery/socket"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/cache"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

func main() {
	conf := configs.InitConfig()
	log := logger.NewLogger(logger.TRACE)
	ch := configs.InitCache(conf)

	const ttl = 24 * time.Hour

	userRepo := cache.NewUserCH(ch, ttl)
	roomRepo := cache.NewRoomCh(ch, ttl)

	userService := service.NewUserService(userRepo, log.WithComponent("UserService"))
	_ = service.NewRoomService(roomRepo, userRepo, log.WithComponent("RoomService"))

	wsHandler := socket.NewWSHandler(
		log.WithComponent("WSHandler"),
		userService,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler.WSHandlerS)
	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	log.Infof("server listening on :%s", conf.ServerPort)
	if err := http.ListenAndServe(":"+conf.ServerPort, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
