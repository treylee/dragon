package main

import (
	"gdragon/internal/router"
	"gdragon/internal/websocket"
	"github.com/sirupsen/logrus" 
	_ "net/http/pprof"
	"net/http"
	"runtime"

)

func main() {

	r := router.SetupRouter()

	go websocket.StartServer("3001")

	logrus.Infof("Starting server on port 8080")
	logrus.Infof("Number of Go threads (goroutines): %d", runtime.NumGoroutine())
	logrus.Infof("Number of CPU cores being utilized by Go runtime: %d", runtime.GOMAXPROCS(0))
	logrus.Infof("Number of available CPU cores: %d", runtime.NumCPU())
	go func() {
		logrus.Info("Starting pprof server on :6060")
		err := http.ListenAndServe(":6060", nil) // Default pprof handler is registered by net/http/pprof
		if err != nil {
			logrus.Fatalf("Error starting pprof server: %v", err)
		}
	}()
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}


}
