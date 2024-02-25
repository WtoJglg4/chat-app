package main

import (
	"log"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../index.html")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(r.Header["Connection"])
		log.Panicln("upgrade: ", err)
		return
	}
	defer ws.Close()

	log.Println("client successfully connected...")
	log.Printf("client: %s", read(ws))
}

func read(ws *ws.Conn) []byte {

	_, byteMsg, err := ws.ReadMessage()
	if err != nil {
		log.Println("read websocket:", err)
		return nil
	}

	return byteMsg
}

func setUpRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	server := http.Server{
		Addr:              "127.0.0.1:8080",
		ReadHeaderTimeout: 3 * time.Second,
	}
	setUpRoutes()
	log.Println("server starts...")
	log.Fatal(server.ListenAndServe())
}
