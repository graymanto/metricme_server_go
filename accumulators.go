package main

import (
	"math"
	"sort"
)

type gaugeValues struct {
	values map[string]int
	syncID int
}

type timerValues struct {
	values        map[string][]int
	flushInterval int
	syncID        int
}

type counterValues struct {
	values        map[string]float32
	flushInterval int
	syncID        int
}

type setValues struct {
	values map[string]map[string]bool
	syncID int
}

type statsToFlush struct {
	syncID   int
	statType rune
	values   interface{}
}

type accCounterMetrics struct {
	values map[string]float32
	rates  map[string]float32
}

type timerMetrics struct {
	std        float64
	upper      int
	lower      int
	count      int
	countPs    float64
	sum        int
	sumSquares int
	mean       float64
	median     float64
}

// TODO: get rid of this type, just send values?
type accTimerMetrics struct {
	values map[string]*timerMetrics
}

var flusher = make(chan statsToFlush)

var accGauge = make(chan gaugeValues)
var accTimer = make(chan timerValues)
var accCounter = make(chan counterValues)
var accSet = make(chan setValues)

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func standardDeviation(series []int, mean float64, count int) float64 {
	var sumOfDiffs float64

	for _, td := range series {
		sumOfDiffs += (float64(td) - mean) * (float64(td) - mean)
	}

	variance := sumOfDiffs / float64(count)
	return round(math.Sqrt(variance), 2)
}

func calculateStatsForSeries(series []int, flushInterval int) *timerMetrics {
	sort.Ints(series)

	currentData := &timerMetrics{}
	count := len(series)
	countPs := float64(count) / (float64(flushInterval) / 1000)

	min := series[0]
	max := series[count-1]
	sum := 0
	sumSquares := 0

	for _, td := range series {
		sum += td
		sumSquares += td * td
	}

	mean := float64(sum) / float64(count)

	var mid int = int(math.Floor(float64(count) / 2))

	var median float64
	if count%2 > 0 || count < 2 {
		median = float64(series[mid])
	} else {
		median = (float64(series[mid-1]) + float64(series[mid])) / 2
	}

	stddev := standardDeviation(series, mean, count)

	currentData.mean = mean
	currentData.lower = min
	currentData.upper = max
	currentData.count = count
	currentData.countPs = countPs
	currentData.sum = sum
	currentData.sumSquares = sumSquares
	currentData.std = stddev
	currentData.median = median

	return currentData
}

func makeTimerMetrics(t timerValues) accTimerMetrics {
	allMetrics := accTimerMetrics{make(map[string]*timerMetrics)}

	for k, val := range t.values {
		count := len(val)

		if count == 0 {
			// TODO: is this a sensible default?
			allMetrics.values[k] = &timerMetrics{}
			continue
		}

		currentData := calculateStatsForSeries(val, t.flushInterval)

		allMetrics.values[k] = currentData
	}

	return allMetrics
}

func gaugeAccum() {
	for v := range accGauge {
		flusher <- statsToFlush{v.syncID, 'g', v.values}
	}
}

func timerAccum() {
	for v := range accTimer {
		metrics := makeTimerMetrics(v)
		flusher <- statsToFlush{v.syncID, 't', metrics}
	}
}

func counterAccum() {
	for v := range accCounter {
		rates := make(map[string]float32)

		for k, val := range v.values {
			rates[k] = val / (float32(v.flushInterval) / 1000)
		}
		acc := accCounterMetrics{v.values, rates}
		flusher <- statsToFlush{v.syncID, 'c', acc}
	}
}

func setAccum() {
	for v := range accSet {
		flusher <- statsToFlush{v.syncID, 's', v.values}
	}
}

func startAccumulators() {
	go gaugeAccum()
	go timerAccum()
	go counterAccum()
	go setAccum()
}
