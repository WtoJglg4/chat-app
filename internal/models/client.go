package models

import (
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
)

type Client struct {
	Hub  *Hub
	Conn *ws.Conn
	Send chan []byte
	Name string
}

// conn -> hub в goroutine
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//тут должен быть анмаршал
		//далее броадкаст хабом структуры Message (а не []byte) по всем клиентам
		//хотя может и не надо(есть смысл анмаршалить на стороне клиента при выводе сообщения на экран)
		// msg := Message{}
		// err := c.Conn.ReadJSON(&msg)

		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error2: %v", err)
			break
		}
		c.Hub.SendAll(msg)
	}
}

// hub -> conn
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(ws.TextMessage)
			if err != nil {
				log.Print("client: cannot create conn writer\n")
				return
			}
			w.Write(msg)

			for i := 0; i < len(c.Send); i++ {
				w.Write(newLine)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
