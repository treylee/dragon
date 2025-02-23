package main

import (
	"gdragon/internal/router"
	"gdragon/internal/websocket"
	"log"
)

func main() {
	r := router.SetupRouter()

	go websocket.StartServer("3001")

	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
