package runner

import (
	"gdragon/internal/metrics"
	"gdragon/internal/websocket"
	"gdragon/internal/utils"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	
)

type TestRunner struct {
	running           bool
	metrics           *metrics.TestMetrics
	wg                sync.WaitGroup
	requestsPerSecond int
	duration          time.Duration
	threadCounter     int
}

func NewTestRunner(requestsPerSecond int, duration time.Duration) *TestRunner {
	return &TestRunner{
		running:           false,
		metrics:           &metrics.TestMetrics{},
		requestsPerSecond: requestsPerSecond,
		duration:          duration,
		threadCounter:     0,
	}
}

func (r *TestRunner) StartTest() {
	logrus.Infof("Starting Test")
	r.running = true
	r.metrics = &metrics.TestMetrics{} // Reset metrics

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	endTime := time.Now().Add(r.duration)
	logrus.Infof("Test will run until: %v (duration: %v)", endTime, r.duration)

	for time.Now().Before(endTime) && r.running {

		<-ticker.C // Blocks until tick fires

		logrus.Infof("Tick - Spawning %d threads", r.requestsPerSecond)

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
				//logrus.Infof("Starting thread #%d", threadID)
				r.runTest()
			}(i + 1)// start the count from thread 1 insteadof0
		}
	}

	r.running = false
	r.wg.Wait()

	r.calculateFinalMetrics()
	websocket.BroadcastMetrics(r.metrics)

}
func (r *TestRunner) calculateFinalMetrics(){
		// Calculate metrics after test
		duration := time.Second * 10
		r.metrics.RequestPerSecond = r.metrics.Requests / int(duration.Seconds())
		r.metrics.AvgResponseTime = float64(r.metrics.ResponseTime) / float64(r.metrics.Requests)
	
		sort.Ints(r.metrics.ResponseTimes)
		r.metrics.P50ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes,50)
		r.metrics.P95ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes,95)
		r.metrics.P99ResponseTime = utils.CalculatePercentile(r.metrics.ResponseTimes,99)
	
		totalRequests := r.metrics.Requests + r.metrics.FailedRequests
		if totalRequests > 0 {
			r.metrics.ErrorRate = (float64(r.metrics.FailedRequests) / float64(totalRequests)) * 100
		}
	
		logrus.Infof("Test completed. Total Requests: %d, Failed Requests: %d", r.metrics.Requests, r.metrics.FailedRequests)
}

func (r *TestRunner) runTest() {
	start := time.Now()
	time.Sleep(time.Millisecond * 200) // Simulated response time

	responseTime := int(time.Since(start).Milliseconds())

	r.metrics.Lock()
	r.metrics.Requests++
	r.metrics.ResponseTimes = append(r.metrics.ResponseTimes, responseTime)
	r.metrics.ResponseTime += responseTime // Update the ResponseTime field each time a response is recorded
	r.metrics.Unlock()

	// Send live metrics update
	websocket.BroadcastMetrics(r.metrics)
}

func (r *TestRunner) IsRunning() bool {
	return r.running
}

func (r *TestRunner) GetMetrics() *metrics.TestMetrics {
	return r.metrics
}
