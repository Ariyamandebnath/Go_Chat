package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		fmt.Println("Received:", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected:", r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return // Return to avoid calling reader on a nil ws
	}

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Simple Server")
	})

	http.HandleFunc("/ws", serveWs)
}

func main() {
	setupRoutes()
	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
