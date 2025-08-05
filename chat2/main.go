package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var mapWsConn = make(map[string]*websocket.Conn)

func main() {
	http.HandleFunc("/chat", LoadPageChat)
	http.HandleFunc("/ws", InitWebsocket)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func LoadPageChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(w, "%s", "error")
		return
	}
	fmt.Println(path)

	content, err := os.ReadFile(path + "\\chat.html")
	if err != nil {
		fmt.Fprintf(w, "%s", "error")
		return
	}

	fmt.Fprintf(w, "%s", content)
}

func InitWebsocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	channel := r.URL.Query().Get("channel")
	if channel == "" {
		http.Error(w, "Missing channel parameter", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Invalid origin", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer func() {
		conn.Close()
		delete(mapWsConn, channel) // Clean up the connection on disconnect
		fmt.Printf("Connection closed for channel: %s\n", channel)
	}()

	mapWsConn[channel] = conn
	fmt.Printf("New connection established for channel: %s\n", channel)

	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}
		fmt.Printf("Received from channel %s: %v\n", channel, msg)

		otherConn := getConn(channel)
		if otherConn == nil {
			fmt.Println("No other connection found for channel:", channel)
			continue
		}

		err = otherConn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error writing JSON to other connection:", err)
			break
		}
	}
}

func getConn(channel string) *websocket.Conn {
	for key, conn := range mapWsConn {
		if key != channel && conn != nil {
			return conn
		}
	}
	return nil
}
