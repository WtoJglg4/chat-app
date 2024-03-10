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

	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s\n", err.Error())
	}

	hub := models.NewHub()
	go hub.Run()

	service := service.NewService()
	handler := handlers.NewHandler(service, hub)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Printf("server error: %v\n", err)
		}
	}()
	log.Printf("server started on port: %s\n", viper.GetString("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
