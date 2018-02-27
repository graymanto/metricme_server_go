package main

import "math"

type gaugeValues struct {
	values map[string]int
	syncID int
}

type timerValues struct {
	values map[string][]int
	syncID int
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
	countPs    int
	sum        int
	sumSquares int
	mean       float64
	median     int
	percent    float64
}

// TODO: get rid of this type, just send values?
type accTimerMetrics struct {
	values map[string]timerMetrics
}

var flusher = make(chan statsToFlush)

var accGauge = make(chan gaugeValues)
var accTimer = make(chan timerValues)
var accCounter = make(chan counterValues)
var accSet = make(chan setValues)

func makeTimerMetrics(t timerValues) accTimerMetrics {
	allMetrics := accTimerMetrics{make(map[string]timerMetrics)}

	// TODO: values should be sorted? Need this for median

	for k, val := range t.values {
		count := len(val)
		currentData := timerMetrics{}

		if count == 0 {
			// TODO: is this a sensible default?
			allMetrics.values[k] = currentData
			continue
		}

		min := val[0]
		max := val[count-1]
		sum := 0
		sumSquares := 0

		for td := range val {
			sum += td
			sumSquares += td * td

			// temp until sorted implementation
			if td < min {
				min = td
			}
			if td > max {
				max = td
			}
		}

		mean := float64(sum) / float64(count)

		var sumOfDiffs float64

		for td := range val {
			sumOfDiffs += (float64(td) - mean) * (float64(td) - mean)
		}

		stddev := math.Sqrt(sumOfDiffs / float64(count))

		currentData.mean = mean
		currentData.lower = min
		currentData.upper = max
		currentData.count = count
		currentData.sum = sum
		currentData.sumSquares = sumSquares
		currentData.std = stddev

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
