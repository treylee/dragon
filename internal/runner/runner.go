package runner

import (
    "gdragon/database/local"
    "gdragon/internal/metrics"
    "gdragon/internal/websocket"
    "net/http"
    "sync/atomic"
    "time"

    "github.com/sirupsen/logrus"
)

func (r *TestRunner) StartTest() {
    logrus.Infof("Starting Test: %s", r.testID)
    r.running.Store(true)
    r.metrics = &metrics.TestMetrics{}

    db, err := local.SetupDatabase(r.testID)
    if err != nil {
        logrus.Errorf("Failed to set up database for test %s: %v", r.testID, err)
        return
    }

    defer db.Close()

    workerCount := r.requestsPerSecond
    for i := 0; i < workerCount; i++ {
        go r.worker()
    }

    expectedRequests := r.requestsPerSecond * int(r.duration.Seconds())

    for i := 0; i < expectedRequests && r.running.Load(); i++ {
        r.jobChannel <- struct{}{} 
        time.Sleep(time.Second / time.Duration(r.requestsPerSecond)) 
    }

    close(r.jobChannel)
    r.running.Store(false)
    r.wg.Wait() 

    r.calculateFinalMetrics()
    websocket.BroadcastMetrics(r.testID, r.metrics)
    websocket.NotifyTestCompletion(r.testID)
    local.SaveTestResults(db, r.metrics)
}

func (r *TestRunner) worker() {
    for range r.jobChannel {
        if !r.running.Load() {
            return
        }
        r.wg.Add(1)
        r.fire()
        r.wg.Done()
    }
}

func (r *TestRunner) fire() {
    start := time.Now()
    resp, err := http.Get("https://httpbin.org/get")
    responseTime := int(time.Since(start).Milliseconds())

    atomic.AddInt64(&r.metrics.Requests, 1)
    atomic.AddInt64(&r.metrics.ResponseTime, int64(responseTime))

    if err != nil || resp.StatusCode >= 400 {
        atomic.AddInt64(&r.metrics.FailedRequests, 1)
        logrus.Warnf("Request failed: %v, StatusCode: %d", err, resp.StatusCode)
    }

    if resp != nil {
        resp.Body.Close()
    }

    if atomic.LoadInt64(&r.metrics.Requests)%100 == 0 {
        websocket.BroadcastMetrics(r.testID, r.metrics)
    }
}

func (r *TestRunner) calculateFinalMetrics() {
    totalRequests := atomic.LoadInt64(&r.metrics.Requests)
    failedRequests := atomic.LoadInt64(&r.metrics.FailedRequests)
    totalResponseTime := atomic.LoadInt64(&r.metrics.ResponseTime)

    if totalRequests > 0 {
        r.metrics.AvgResponseTime = float64(totalResponseTime) / float64(totalRequests)
    }

    if totalRequests+failedRequests > 0 {
        r.metrics.ErrorRate = (float64(failedRequests) / float64(totalRequests+failedRequests)) * 100
    }

    r.metrics.TestID = r.testID
    r.metrics.TestName = r.testName
    r.metrics.RequestsPerSecond = r.requestsPerSecond
    r.metrics.TestDuration = int(r.duration.Seconds())

    logrus.Infof("Test Completed: %s | Total Requests: %d | Failed Requests: %d", r.testID, totalRequests, failedRequests)
}

func (r *TestRunner) IsRunning() bool {
    return r.running.Load()
}

func (r *TestRunner) GetMetrics() *metrics.TestMetrics {
    return r.metrics
}
func (r *TestRunner) GetTestID() string {
    return r.testID
}
