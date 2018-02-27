package main

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

var flusher = make(chan statsToFlush)

var accGauge = make(chan gaugeValues)
var accTimer = make(chan timerValues)
var accCounter = make(chan counterValues)
var accSet = make(chan setValues)

func makeTimerMetrics(t timerValues) {

}

func gaugeAccum() {
	for v := range accGauge {
		flusher <- statsToFlush{v.syncID, 'g', v.values}
	}
}

func timerAccum() {
	for v := range accTimer {
		flusher <- statsToFlush{v.syncID, 't', 0}
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
