package runner

import (
	"github.com/sirupsen/logrus"
	"gdragon/internal/metrics"
	"sort"
	"sync"
	"time"
)

type TestRunner struct {
	running bool
	metrics *metrics.TestMetrics
	log     *logrus.Logger
}

func NewTestRunner() *TestRunner {
	// Create a new instance of Logrus logger
	logger := logrus.New()
	// Set the log level
	logger.SetLevel(logrus.InfoLevel)
	// Set the log format (can be changed to JSON if preferred)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return &TestRunner{
		running: false,
		metrics: &metrics.TestMetrics{},
		log:     logger,
	}
}

func (r *TestRunner) StartTest() {
	r.running = true
	startTime := time.Now()

	var wg sync.WaitGroup

	r.log.Info("Starting performance test...")

	for i := 0; i < 150; i++ {
		wg.Add(1)
		go r.runTest(i, &wg)
	}

	wg.Wait()

	duration := time.Since(startTime)
	r.metrics.RequestRate = float64(r.metrics.Requests) / duration.Seconds()
	r.metrics.AvgResponseTime = float64(r.metrics.ResponseTime) / float64(r.metrics.Requests)

	sort.Ints(r.metrics.ResponseTimes)
	r.metrics.P50ResponseTime = r.CalculatePercentile(50)
	r.metrics.P95ResponseTime = r.CalculatePercentile(95)
	r.metrics.P99ResponseTime = r.CalculatePercentile(99)

	if r.metrics.Requests > 0 {
		r.metrics.ErrorRate = (float64(r.metrics.FailedRequests) / float64(r.metrics.Requests)) * 100
	}

	// Log the metrics struct in a cleaner way
	r.log.Infof("\nTest Completed!\nMetrics: %+v", *r.metrics)
}

func (r *TestRunner) runTest(userID int, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Millisecond * 200)

	if userID%2 == 0 {
		r.metrics.FailedRequests++
		r.log.WithFields(logrus.Fields{
			"user_id": userID,
			"status":  "failed",
		}).Info("Request Failed")
	} else {
		r.metrics.Requests++
		r.metrics.ResponseTime += 200
		r.metrics.ResponseTimes = append(r.metrics.ResponseTimes, 200)
		r.log.WithFields(logrus.Fields{
			"user_id":     userID,
			"response_ms": 200,
			"status":      "success",
		}).Info("Request Success")
	}

	if r.metrics.ResponseTime > r.metrics.MaxResponseTime {
		r.metrics.MaxResponseTime = r.metrics.ResponseTime
	}
}

func (r *TestRunner) CalculatePercentile(percentile int) int {
	index := (percentile * len(r.metrics.ResponseTimes)) / 100
	if index >= len(r.metrics.ResponseTimes) {
		return r.metrics.ResponseTimes[len(r.metrics.ResponseTimes)-1]
	}
	return r.metrics.ResponseTimes[index]
}

func (r *TestRunner) GetMetrics() *metrics.TestMetrics {
	return r.metrics
}

func (r *TestRunner) IsRunning() bool {
	return r.running
}
