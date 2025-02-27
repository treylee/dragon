package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
	"gdragon/internal/runner"
	"gdragon/internal/models"
)

var (
	testRunners map[string]*runner.TestRunner
	mu           sync.Mutex
)

func init() {

	testRunners = make(map[string]*runner.TestRunner)
}

func StartTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var testRequest models.StartTestRequest

	if err := c.ShouldBindJSON(&testRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	testID := "id-" + uuid.New().String() + "-" + time.Now().Format("20060102150405")

	if _, exists := testRunners[testID]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Test with this ID already exists", "testID": testID})
		return
	}

	testRunner := runner.NewTestRunner(testRequest.RequestPerSecond, time.Second*time.Duration(testRequest.Duration), testID, testRequest.TestName)

	testRunners[testID] = testRunner

	go func() {
		testRunner.StartTest()
		delete(testRunners, testID) 
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Test started", "testID": testID})
}

func TestStatus(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	// Get testID from URL params
	testID := c.Param("testID")

	testRunner, exists := testRunners[testID]
	if !exists || !testRunner.IsRunning() {
		c.JSON(http.StatusOK, gin.H{"status": "No test running"})
		return
	}

	c.JSON(http.StatusOK, testRunner.GetMetrics())
}

func GetTests(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	// Return the list of all running tests with their IDs
	tests := make([]gin.H, 0, len(testRunners))
	for id, runner := range testRunners {
		tests = append(tests, gin.H{
			"testID": id,
			"status": runner.IsRunning(),
		})
	}

	c.JSON(http.StatusOK, tests)
}
