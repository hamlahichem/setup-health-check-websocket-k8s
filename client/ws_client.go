package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {

	var SERVER = "localhost:8080"
	var PATH = "/ws"
	fmt.Println("Connecting to:", SERVER, "at", PATH)

	URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
	c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, []byte("msg"))
	if err != nil {
		log.Println("Write error:", err)
		return
	} else {

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("ReadMessage() error:", err)
			return
		}
		log.Printf("Received: %s", message)
	}
	defer c.Close()
}
