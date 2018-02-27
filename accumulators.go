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
	values map[string]float32
	syncID int
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

var flusher = make(chan statsToFlush)

var accGauge = make(chan gaugeValues)
var accTimer = make(chan timerValues)
var accCounter = make(chan counterValues)
var accSet = make(chan setValues)

func gaugeAccum() {
	for v := range accGauge {
		flusher <- statsToFlush{v.syncID, 'g', 0}
	}
}

func timerAccum() {
	for v := range accTimer {
		flusher <- statsToFlush{v.syncID, 't', 0}
	}
}

func counterAccum() {
	for v := range accCounter {
		flusher <- statsToFlush{v.syncID, 'c', 0}
	}
}

func setAccum() {
	for v := range accSet {
		flusher <- statsToFlush{v.syncID, 's', 0}
	}
}

func startAccumulators() {
	go gaugeAccum()
	go timerAccum()
	go counterAccum()
	go setAccum()
}
