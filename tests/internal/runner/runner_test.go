package runner

import (
	"gdragon/internal/runner" // Update the import path to point to the correct package
	"sort"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestTestRunner_RealisticRequestSimulation(t *testing.T) {
	testRunner := runner.NewTestRunner()


	for i := 0; i < 150; i++ {
		if i%3 == 0 {
			testRunner.GetMetrics().Lock()
			testRunner.GetMetrics().FailedRequests++
			testRunner.GetMetrics().Unlock()
		} else {
			testRunner.GetMetrics().Lock()
			testRunner.GetMetrics().Requests++
			testRunner.GetMetrics().ResponseTime += 250
			testRunner.GetMetrics().ResponseTimes = append(testRunner.GetMetrics().ResponseTimes, 250)
			testRunner.GetMetrics().Unlock()
		}
	}

	// Calculate request rate and average response time
	duration := time.Second * 10
	testRunner.GetMetrics().RequestRate = float64(testRunner.GetMetrics().Requests) / duration.Seconds()
	testRunner.GetMetrics().AvgResponseTime = float64(testRunner.GetMetrics().ResponseTime) / float64(testRunner.GetMetrics().Requests)

	sort.Ints(testRunner.GetMetrics().ResponseTimes)
	testRunner.GetMetrics().P50ResponseTime = testRunner.CalculatePercentile(50)
	testRunner.GetMetrics().P95ResponseTime = testRunner.CalculatePercentile(95)
	testRunner.GetMetrics().P99ResponseTime = testRunner.CalculatePercentile(99)

	if testRunner.GetMetrics().Requests > 0 {
		testRunner.GetMetrics().ErrorRate = (float64(testRunner.GetMetrics().FailedRequests) / float64(testRunner.GetMetrics().Requests)) * 100
	}

	assert.Equal(t, 100, testRunner.GetMetrics().Requests, "Expected Requests to be 100")
	assert.Equal(t, 50, testRunner.GetMetrics().FailedRequests, "Expected FailedRequests to be 50")
	assert.Equal(t, 50.0, testRunner.GetMetrics().ErrorRate, "Expected ErrorRate to be 50.0")

	expectedPercentile := 250
	assert.Equal(t, expectedPercentile, testRunner.GetMetrics().P50ResponseTime, "Expected P50ResponseTime to be 250")
	assert.Equal(t, expectedPercentile, testRunner.GetMetrics().P95ResponseTime, "Expected P95ResponseTime to be 250")
	assert.Equal(t, expectedPercentile, testRunner.GetMetrics().P99ResponseTime, "Expected P99ResponseTime to be 250")
}
