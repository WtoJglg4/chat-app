package handlers

import (
	"github/chat-app/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

func (h *Handler) chatPage(c *gin.Context) {
	c.File("index.html")
}

func (h *Handler) styles(c *gin.Context) {
	c.File("front/styles.css")
}

func (h *Handler) script(c *gin.Context) {
	c.File("front/script.js")
}

func (h *Handler) favicon(c *gin.Context) {
	c.File("front/favicon.ico")
}

var upgrader = ws.Upgrader{}

func (h *Handler) websocket(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(c.Request.Header["Connection"])
		log.Panicln("upgrade: ", err)
		return
	}

	log.Println("client successfully connected to websocket...")
	client := &models.Client{
		Hub:  h.hub,
		Conn: ws,
		Send: make(chan []byte),
	}
	h.hub.Register(client)

	go client.ReadPump()
	go client.WritePump()
}
