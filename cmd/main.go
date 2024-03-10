package main

import (
	server "github/chat-app"
	"github/chat-app/internal/models"
	handlers "github/chat-app/pkg"
	"github/chat-app/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	hub := models.NewHub()
	go hub.Run()

	service := service.NewService()
	handler := handlers.NewHandler(service, hub)
	srv := new(server.Server)
	go func() {
		if err := srv.Run("8080", handler.InitRoutes()); err != nil {
			log.Printf("server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
