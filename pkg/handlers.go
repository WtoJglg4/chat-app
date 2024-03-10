package handlers

import (
	"github/chat-app/internal/models"
	"github/chat-app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	hub      *models.Hub
}

func NewHandler(service *service.Service, hub *models.Hub) *Handler {
	return &Handler{services: service, hub: hub}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/", h.chatPage)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	chat := router.Group("/chat")
	{
		chat.GET("", h.chatPage)
		chat.GET("/ws", h.websocket)
	}

	front := router.Group("/front")
	{
		front.GET("/styles.css", h.styles)
		front.GET("/script.js", h.script)
		front.GET("/favicon.ico", h.favicon)
	}

	return router
}
