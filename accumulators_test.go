package main

import (
	"reflect"
	"testing"
)

var metricTests = []struct {
	in  []int
	out *timerMetrics
}{
	{[]int{90, 9, 59, 27, 76, 24, 3, 49, 3, 96},
		&timerMetrics{
			33.72, 96, 3, 10, 2.0, 436, 30378, 43.6, 38.0}},
	{[]int{13, 8, 30, 57, 55, 69, 10, 48, 81, 82},
		&timerMetrics{27.09, 82, 8, 10, 2.0, 453, 27857, 45.3, 51.5}},
	{[]int{14, 40, 56, 70, 87, 44, 87, 57, 0, 25},
		&timerMetrics{27.82, 87, 0, 10, 2.0, 480, 30780, 48.0, 50.0}},
	{[]int{15, 13, 42, 5, 67, 91, 93, 19, 22, 11},
		&timerMetrics{32.06, 93, 5, 10, 2.0, 378, 24568, 37.8, 20.5}},
	{[]int{80, 45, 89, 38, 3, 24, 65, 44, 72, 27},
		&timerMetrics{25.95, 89, 3, 10, 2.0, 487, 30449, 48.7, 44.5}},
	{[]int{5},
		&timerMetrics{0.0, 5, 5, 1, 0.2, 5, 25, 5.0, 5.0}},
}

// TestCalculateStatsForSeries tests that stats are calculated correctly for a series
// of integers.
func TestCalculateStatsForSeries(t *testing.T) {
	const flushInterval = 5000

	for _, tt := range metricTests {
		stats := calculateStatsForSeries(tt.in, flushInterval)
		if !reflect.DeepEqual(stats, tt.out) {
			t.Error("Parse result not as expected", stats, tt.out)
		}
	}
}

// TestMakeTimerMetrics essentially runs the same tests as TestCalculateStatsForSeries
// but puts the inputs in the formed of a keyed map and passes them in to makeTimerMetrics,
// then checks the results are correctly placed into a map again.
func TestMakeTimerMetrics(t *testing.T) {
	var testInput = timerValues{}
	testInput.flushInterval = 5000
	testInput.values = make(map[string][]int)
	const keyBase = "val"

	for i, tt := range metricTests {
		testInput.values[keyBase+string(i)] = tt.in
	}

	result := makeTimerMetrics(testInput)

	for i, tt := range metricTests {
		key := keyBase + string(i)
		if !reflect.DeepEqual(result.values[key], tt.out) {
			t.Error("Parse result not as expected", result.values[key], tt.out)
		}
	}
}
