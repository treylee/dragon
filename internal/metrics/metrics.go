package metrics

import "sync"

type TestMetrics struct {
	mu sync.Mutex
	Requests        int
	ResponseTime    int
	AvgResponseTime float64
	ErrorRate       float64
	RequestRate     float64
	ConcurrentUsers int
	FailedRequests  int
	ResponseTimes   []int
	MaxResponseTime int
	P50ResponseTime int
	P95ResponseTime int
	P99ResponseTime int
	CpuUsage        float64
	MemoryUsage     float64
}

func (m *TestMetrics) Lock() {
	m.mu.Lock()
}

func (m *TestMetrics) Unlock() {
	m.mu.Unlock()
}

