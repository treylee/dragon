package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"gdragon/internal/runner"
)

var (
	testRunner *runner.TestRunner
	mu         sync.Mutex
)

// init a perftest
func StartTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if testRunner == nil {
		testRunner = runner.NewTestRunner()
	}

	if testRunner.IsRunning() {
		c.JSON(http.StatusConflict, gin.H{"message": "Test is already running"})
		return
	}

	go testRunner.StartTest()

	c.JSON(http.StatusOK, gin.H{"message": "Test started"})
}

// return the status//current metrics-
func TestStatus(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if testRunner == nil || !testRunner.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"status": "No test running"})
		return
	}

	c.JSON(http.StatusOK, testRunner.GetMetrics())
}
