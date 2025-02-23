package utils

import (
	"sort"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestUtils_CalculatePercentiles(t *testing.T) {

	testCases := []struct {
        responseTimes []int
        percentile    int
        expected      int
    }{
        {
            responseTimes: []int{100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000},
            percentile:    50,
            expected:      500,
        },
        {
            responseTimes: []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
            percentile:    50,
            expected:      50,
        },
        {
            responseTimes: []int{300, 400, 500, 600, 700, 800, 900, 1000},
            percentile:    95,
            expected:      1000,
        },
        {
            responseTimes: []int{5, 15, 25, 35, 45, 55},
            percentile:    99,
            expected:      55,
        },
    }

    for _, tc := range testCases {
        sort.Ints(tc.responseTimes)

        result := CalculatePercentile(tc.responseTimes, tc.percentile)

        assert.Equal(t, tc.expected, result, "Expected percentile %d to be %d", tc.percentile, tc.expected)
    }
}

