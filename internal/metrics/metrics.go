package metrics

import (
	"sync"
	"sync/atomic"
)

type TestMetrics struct {
	mu                sync.Mutex
	Requests          int64   
	RequestsPerSecond int     
	AvgResponseTime   float64
	ErrorRate         float64
	ConcurrentUsers   int
	FailedRequests    int64 
	ResponseTime      int64 
	ResponseTimes     []int
	MaxResponseTime   int
	P50ResponseTime   int
	P95ResponseTime   int
	P99ResponseTime   int
	CpuUsage          float64
	MemoryUsage       float64
	Url               string
	TestCompleted     bool
	TestDuration      int
	TestName          string
	TestID            string
	CreatedAt         string
}

func (m *TestMetrics) Lock() {
	m.mu.Lock()
}

func (m *TestMetrics) Unlock() {
	m.mu.Unlock()
}

// Atomic Add method for Requests
func (m *TestMetrics) IncrementRequests() {
	atomic.AddInt64(&m.Requests, 1)
}

// Atomic Add method for FailedRequests
func (m *TestMetrics) IncrementFailedRequests() {
	atomic.AddInt64(&m.FailedRequests, 1)
}

// Atomic Add method for ResponseTime
func (m *TestMetrics) AddResponseTime(rt int64) {
	atomic.AddInt64(&m.ResponseTime, rt)
}

// Atomic Get method for Requests
func (m *TestMetrics) GetRequests() int64 {
	return atomic.LoadInt64(&m.Requests)
}

// Atomic Get method for FailedRequests
func (m *TestMetrics) GetFailedRequests() int64 {
	return atomic.LoadInt64(&m.FailedRequests)
}

// Atomic Get method for ResponseTime
func (m *TestMetrics) GetResponseTime() int64 {
	return atomic.LoadInt64(&m.ResponseTime)
}
