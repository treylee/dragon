package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
	"gdragon/internal/runner"
	"github.com/sirupsen/logrus"

)

var (
	testRunners map[string]*runner.TestRunner
	mu           sync.Mutex
)

func init() {

	testRunners = make(map[string]*runner.TestRunner)
}
var log = logrus.New()

func StartTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var testRequest StartTestRequest

	if err := c.ShouldBindJSON(&testRequest); err != nil {
		log.WithError(err).Error("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	log.WithFields(logrus.Fields{
		"TestName":        testRequest.TestName,
		"RequestsPerSecond": testRequest.RequestPerSecond,
		"Duration":        testRequest.Duration,
		"Url":				testRequest.Url,
	}).Info("Received Test Request")

	testID := "id-" + uuid.New().String() + "-" + time.Now().Format("20060102150405")

	if _, exists := testRunners[testID]; exists {
		log.WithFields(logrus.Fields{
			"testID": testID,
		}).Warn("Test with this ID already exists")

		c.JSON(http.StatusConflict, gin.H{"message": "Test with this ID already exists", "testID": testID})
		return
	}

	startTime := time.Now()

	testRunner := runner.NewTestRunner(testRequest.RequestPerSecond, time.Second*time.Duration(testRequest.Duration), testID, testRequest.TestName)

	testRunners[testID] = testRunner

	go func() {
		testRunner.StartTest()

		elapsedTime := time.Since(startTime)
		log.WithFields(logrus.Fields{
			"testID":     testID,
			"elapsedTime": elapsedTime,
		}).Info("Test completed")
		delete(testRunners, testID) 
	}()

	log.WithFields(logrus.Fields{
		"testID": testID,
	}).Info("Test started successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Test started", "testID": testID})
}


func TestStatus(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

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

	tests := make([]gin.H, 0, len(testRunners))
	for id, runner := range testRunners {
		tests = append(tests, gin.H{
			"testID": id,
			"status": runner.IsRunning(),
		})
	}

	c.JSON(http.StatusOK, tests)
}
