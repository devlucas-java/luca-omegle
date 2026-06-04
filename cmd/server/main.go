package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
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
	roomRepo := cache.NewRoomCh(ch, ttl, "room:")

	userService := service.NewUserService(userRepo, log.WithComponent("UserService"))
	roomService := service.NewRoomService(roomRepo, userRepo, log.WithComponent("RoomService"))

	hub := socket.NewHub(
		log.WithComponent("Hub"),
		roomService,
		userService,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go hub.Run(ctx)

	wsHandler := socket.NewWSHandler(
		log.WithComponent("WSHandler"),
		userService,
		hub,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler.ServeWS)

	srv := &http.Server{
		Addr:    ":" + conf.ServerPort,
		Handler: mux,
	}

	go func() {
		log.Infof("server listening on :%s", conf.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Errorf("graceful shutdown: %v", err)
	}
}
