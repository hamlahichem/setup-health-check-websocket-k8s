package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	// Define a WebSocket upgrade handler
	upgrader := websocket.Upgrader{}
	//started := time.Now()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		// Echo incoming messages back to the client
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("ws server Received message:", string(p))

			err = conn.WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	/*
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			duration := time.Now().Sub(started)
			fmt.Println("healthy is hitted")

			if duration.Seconds() > 60 {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
				fmt.Println("error:  second > 60")
				fmt.Println("healthy status = failed", duration.Seconds())
			} else {
				w.WriteHeader(200)
				w.Write([]byte("ok"))
				fmt.Println("healthy status = ok", duration.Seconds())
			}

		})
	*/

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		var SERVER = "localhost:8080"
		var PATH = "/ws"
		fmt.Println("Connecting to:", SERVER, "at", PATH)

		URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
		c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(500)
			return
		}
		err = c.WriteMessage(websocket.TextMessage, []byte("health-check"))
		if err != nil {
			log.Println("Write error:", err)
			w.WriteHeader(500)
			return
		} else {

			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ReadMessage() error:", err)
				w.WriteHeader(500)
				return
			} else {
				log.Printf("ws client Received: %s", message)
				fmt.Println("health passed: websocker server works")
				w.WriteHeader(200)
				fmt.Fprintf(w, "Oops, the page you requested doesn't exist.")
				if string(message) == "ping" {
					// If the client sends 'ping', respond with 'pong'
					err = c.WriteMessage(websocket.TextMessage, []byte("pong"))
					if err != nil {
						log.Println(err)
						return
					}
				} else if string(message) == "pong" {
					// If the client sends 'pong', you can handle it as desired
					log.Println("Received pong from client")
				} else {
					// Handle other messages from the client
					log.Println("Received:", string(message))
				}

			}

		}

		defer c.Close()
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Oops, the page you requested doesn't exist.")
		fmt.Println(" root is hitted")
	})
	// Start the HTTP server
	log.Println("Starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
