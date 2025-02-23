package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var (
	clients = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			logrus.Infof("WebSocket connection attempt from: %s", r.RemoteAddr)
			return true
		},
	}
	mu sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("WebSocket Upgrade Error: %v", err)
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	logrus.Infof("New WebSocket connection established: %v", conn.RemoteAddr())
	go listenForMessages(conn)
	select {}
}

func listenForMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Warnf("Error reading message from %v: %v", conn.RemoteAddr(), err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
			break
		}
		logrus.Infof("Received message from %v: %s", conn.RemoteAddr(), string(msg))
	}
}

func BroadcastMetrics(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Error marshalling metrics data: %v", err)
		return
	}

	mu.Lock()
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			logrus.Errorf("WebSocket Error: %v, closing connection", err)
			client.Close()
			delete(clients, client)
		}
	}
	mu.Unlock()
	logrus.Infof("Broadcasted metrics to %d clients", len(clients))
}

func StartServer(port string) {
	http.HandleFunc("/ws", HandleConnections)
	logrus.Infof("WebSocket server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatalf("Error starting WebSocket server: %v", err)
	}
}
