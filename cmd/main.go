package main

import (
	"fmt"
	"github/chat-app/internal/models"
	"log"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if r.URL.Path != "/" {
		http.Error(w, fmt.Sprintf("Not Found\n%d\n", http.StatusNotFound), http.StatusNotFound)
		log.Println("http : Not Found ", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Method Not Allowed\n%d\n", http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Println("http : Method Not Allowed ", http.StatusMethodNotAllowed)
		return
	}

	log.Println("serving html...")
	http.ServeFile(w, r, "front/index.html")
}

func styles(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	http.ServeFile(w, r, "front/styles.css")
}

func script(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	http.ServeFile(w, r, "front/script.js")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request, hub *models.Hub) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(r.Header["Connection"])
		log.Panicln("upgrade: ", err)
		return
	}

	log.Println("client successfully connected...")
	client := &models.Client{
		Hub:  hub,
		Conn: ws,
		Send: make(chan []byte),
	}
	hub.Register(client)

	go client.ReadPump()
	go client.WritePump()
}

func setUpRoutes(hub *models.Hub) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/front/styles.css", styles)
	http.HandleFunc("/front/script.js", script)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(w, r, hub)
	})
}

func main() {
	hub := models.NewHub()
	go hub.Run()

	serverAddr := "localhost:8080"

	server := http.Server{
		Addr:              serverAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	setUpRoutes(hub)
	log.Println("server starts on " + serverAddr)
	log.Fatal(server.ListenAndServe())
}
