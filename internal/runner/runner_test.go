package runner

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gdragon/internal/utils"
)

func TestRunner_RealisticRequestSimulation(t *testing.T) {
	testRunner := NewTestRunner(10, time.Second*10, "testID","testName")

	// Simulating 150 requests
	for i := 0; i < 150; i++ {
		testRunner.metrics.Lock()
		if i%3 == 0 {
			testRunner.metrics.FailedRequests++
		} else {
			testRunner.metrics.Requests++
			testRunner.metrics.ResponseTime += 250
			testRunner.metrics.ResponseTimes = append(testRunner.metrics.ResponseTimes, 250)
		}
		testRunner.metrics.Unlock()
	}

	testRunner.metrics.Lock()
	duration := time.Second * 10
	testRunner.metrics.RequestPerSecond = 	testRunner.metrics.Requests / int(duration.Seconds())
	testRunner.metrics.AvgResponseTime = float64(testRunner.metrics.ResponseTime) / float64(testRunner.metrics.Requests)
	
	sort.Ints(testRunner.metrics.ResponseTimes)
	testRunner.metrics.P50ResponseTime = utils.CalculatePercentile(testRunner.metrics.ResponseTimes, 50)
	testRunner.metrics.P95ResponseTime = utils.CalculatePercentile(testRunner.metrics.ResponseTimes, 95)
	testRunner.metrics.P99ResponseTime = utils.CalculatePercentile(testRunner.metrics.ResponseTimes, 99)

	totalRequests := testRunner.metrics.Requests + testRunner.metrics.FailedRequests
	testRunner.metrics.ErrorRate = (float64(testRunner.metrics.FailedRequests) / float64(totalRequests)) * 100
	
	testRunner.metrics.Unlock()

	assert.Equal(t, 100, testRunner.metrics.Requests, "Expected Requests to be 100")
	assert.Equal(t, 50, testRunner.metrics.FailedRequests, "Expected FailedRequests to be 50")
	assert.InDelta(t, 33.33, testRunner.metrics.ErrorRate, 0.1, "Expected ErrorRate to be 33.33%")
	assert.InDelta(t, 10, testRunner.metrics.RequestPerSecond, 0.1, "Expected RequestsPerSecond to be ~10")

	expectedPercentile := 250
	assert.Equal(t, expectedPercentile, testRunner.metrics.P50ResponseTime, "Expected P50ResponseTime to be 250")
	assert.Equal(t, expectedPercentile, testRunner.metrics.P95ResponseTime, "Expected P95ResponseTime to be 250")
	assert.Equal(t, expectedPercentile, testRunner.metrics.P99ResponseTime, "Expected P99ResponseTime to be 250")

}
