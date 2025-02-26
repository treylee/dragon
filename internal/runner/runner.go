package runner

import (
	"gdragon/internal/metrics"
	"gdragon/internal/websocket"
	"gdragon/internal/utils"
	"sort"
	"sync"
	"time"
	"net/http"

	"github.com/sirupsen/logrus"
)

type TestRunner struct {
	running           bool
	metrics           *metrics.TestMetrics
	wg                sync.WaitGroup
	requestsPerSecond int
	duration          time.Duration
	threadCounter     int
	testID            string 
}

func NewTestRunner(requestsPerSecond int, duration time.Duration, testID string) *TestRunner {
	return &TestRunner{
		running:           false,
		metrics:           &metrics.TestMetrics{},
		requestsPerSecond: requestsPerSecond,
		duration:          duration,
		threadCounter:     0,
		testID:            testID, 
	}
}
func (r *TestRunner) GetTestID() string {
	return r.testID
}

func (r *TestRunner) StartTest() {
	logrus.Infof("Starting Test with testID: %s", r.testID)
	r.running = true
	r.metrics = &metrics.TestMetrics{} 

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	endTime := time.Now().Add(r.duration)
	logrus.Infof("Test will run until: %v (duration: %v)", endTime, r.duration)

	for time.Now().Before(endTime) && r.running {
		<-ticker.C // Blocks until tick fires

		logrus.Infof("Tick - Spawning %d threads for testID: %s", r.requestsPerSecond, r.testID)

		for i := 0; i < r.requestsPerSecond; i++ {
			r.wg.Add(1)
			go func(threadID int) {
				defer r.wg.Done()
				defer func() {
					if err := recover(); err != nil {
						logrus.Errorf("Thread #%d panicked: %v", threadID, err)
					}
				}()

				r.threadCounter++
				r.runTest()
			}(i + 1) // Start the count from thread 1 instead of 0
		}
	}

	r.running = false
	r.wg.Wait()

	r.calculateFinalMetrics()
	websocket.BroadcastMetrics(r.testID, r.metrics)
	websocket.NotifyTestCompletion(r.testID)
}

func (r *TestRunner) calculateFinalMetrics() {

	duration := time.Second * 10
	r.metrics.RequestPerSecond = r.metrics.Requests / int(duration.Seconds())
	r.metrics.AvgResponseTime = float64(r.metrics.ResponseTime) / float64(r.metrics.Requests)

	sort.Ints(r.metrics.ResponseTimes)
	r.metrics.P50ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes, 50)
	r.metrics.P95ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes, 95)
	r.metrics.P99ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes, 99)
	
	totalRequests := r.metrics.Requests + r.metrics.FailedRequests
	if totalRequests > 0 {
		r.metrics.ErrorRate = (float64(r.metrics.FailedRequests) / float64(totalRequests)) * 100
	}

	logrus.Infof("Test completed for testID: %s. Total Requests: %d, Failed Requests: %d", r.testID, r.metrics.Requests, r.metrics.FailedRequests)
}

func (r *TestRunner) runTest() {
	start := time.Now()
	resp,err := http.Get("https://httpbin.org/get")
	responseTime := int(time.Since(start).Milliseconds())

	r.metrics.Lock()
	r.metrics.Requests++
	r.metrics.ResponseTimes = append(r.metrics.ResponseTimes, responseTime)
	r.metrics.ResponseTime += responseTime // Update the ResponseTime field each time a response is recorded
	r.metrics.Unlock()

	if err != nil || resp.StatusCode >= 400 {
		r.metrics.FailedRequests++
		logrus.Warnf("Request failed: %v", err)
	}

	if resp != nil {
		resp.Body.Close()
	}

	websocket.BroadcastMetrics(r.testID, r.metrics) 
}

func (r *TestRunner) IsRunning() bool {
	return r.running
}

func (r *TestRunner) GetMetrics() *metrics.TestMetrics {
	return r.metrics
}
