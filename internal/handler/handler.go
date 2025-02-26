package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
	"gdragon/internal/runner"
)

// Global testRunner instance
var (
	testRunner *runner.TestRunner
	mu         sync.Mutex
)

func StartTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	testID := "id-" + uuid.New().String() + "-" + time.Now().Format("20060102150405")

	if testRunner == nil || !testRunner.IsRunning() {
		testRunner = runner.NewTestRunner(10, time.Second*10, testID)
	} else {
		c.JSON(http.StatusConflict, gin.H{"message": "Test is already running", "testID": testRunner.GetTestID()})
		return
	}

	go testRunner.StartTest()

	c.JSON(http.StatusOK, gin.H{"message": "Test started", "testID": testID})
}

func TestStatus(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if testRunner == nil || !testRunner.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"status": "No test running"})
		return
	}

	c.JSON(http.StatusOK, testRunner.GetMetrics())
}

