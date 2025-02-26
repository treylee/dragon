package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var (
	clients = make(map[string]map[*websocket.Conn]bool) // Clients mapped by testID
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true 
		},
	}
	mu sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	testID := r.URL.Query().Get("testid") // Changed to match query parameter
	if testID == "" {
		http.Error(w, "testID query parameter is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("WebSocket Upgrade Error: %v", err)
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	if clients[testID] == nil {
		clients[testID] = make(map[*websocket.Conn]bool)
	}
	clients[testID][conn] = true
	mu.Unlock()

	logrus.Infof("WebSocket connection established for testID: %s", testID)

	go listenForMessages(conn, testID)
}

func listenForMessages(conn *websocket.Conn, testID string) {
	defer func() {
		mu.Lock()
		delete(clients[testID], conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("Error reading message: %v", err)
			return
		}
	}
}

func BroadcastMetrics(testID string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Error marshalling metrics data: %v", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Send the message to all connected clients
	for client := range clients[testID] {
		if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			logrus.Errorf("Error sending message to client: %v", err)
			client.Close()
			delete(clients[testID], client)
		}
	}

	logrus.Infof("Broadcasted metrics for testID %s to %d clients", testID, len(clients[testID]))
}

// NotifyTestCompletion sends a final message indicating test completion and closes all connections
func NotifyTestCompletion(testID string) {
	completionMessage := map[string]interface{}{
		"TestCompleted": true,
	}

	jsonData, err := json.Marshal(completionMessage)
	if err != nil {
		logrus.Errorf("Error marshalling completion message: %v", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for client := range clients[testID] {
		if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			logrus.Errorf("Error sending completion message: %v", err)
		}
		client.Close()
		delete(clients[testID], client)
	}

	logrus.Infof("Notified test completion for testID %s and closed all connections", testID)
}

func StartServer(port string) {
	http.HandleFunc("/ws", HandleConnections)
	logrus.Infof("WebSocket server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatalf("Error starting WebSocket server: %v", err)
	}
}
