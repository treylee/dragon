package runner

import (
	"gdragon/internal/metrics"
	"sync"
	"sync/atomic"
	"time"
)
type TestRunner struct {
    running           atomic.Bool
    metrics           *metrics.TestMetrics
    wg                sync.WaitGroup
    requestsPerSecond int
    duration          time.Duration
    testID            string
    testName          string
    url               string
    jobChannel        chan struct{} 
}

func NewTestRunner(requestsPerSecond int, duration time.Duration, testID, testName string,url string) *TestRunner {
    return &TestRunner{
        metrics:           &metrics.TestMetrics{},
        requestsPerSecond: requestsPerSecond,
        duration:          duration,
        testID:            testID,
        testName:          testName,
        url:               url,
        jobChannel:        make(chan struct{}, requestsPerSecond), 
    }
}
