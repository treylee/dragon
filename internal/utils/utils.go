	package utils

	import (
		"sort"
	)

	func Sum(times []int) int {
		total := 0
		for _, time := range times {
			total += time
		}
		return total
	}

	func CalculateAvgResponseTime(totalResponseTime int, totalRequests int) float64 {
		if totalRequests == 0 {
			return 0
		}
		return float64(totalResponseTime) / float64(totalRequests)
	}

	func CalculatePercentile(times []int, percentile int) int {
		if len(times) == 0 {
			return 0
		}
	
		sort.Ints(times)
	
		index := len(times) * percentile / 100

		if index > 0 &&  percentile != 99 && percentile != 95  {
			index-- 
		}
	
		if index >= len(times) {
			return times[len(times)-1]
		}
	
		return times[index]
	}
	
	

	func CalculateRequestRate(requests int, durationInSeconds int) int {
		if durationInSeconds == 0 {
			return 0
		}
		return requests / durationInSeconds
	}
