package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
	"gdragon/internal/runner"
)

var (
	testRunner *runner.TestRunner
	mu         sync.Mutex
)

func StartTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if testRunner == nil {
		testRunner = runner.NewTestRunner(10, time.Second*10) 
	}

	if testRunner.IsRunning() {
		c.JSON(http.StatusConflict, gin.H{"message": "Test is already running"})
		return
	}

	go testRunner.StartTest()

	c.JSON(http.StatusOK, gin.H{"message": "Test started"})
}

func TestStatus(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, testRunner.GetMetrics())

	if testRunner == nil || !testRunner.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"status": "No test running"})
		return
	}
}
