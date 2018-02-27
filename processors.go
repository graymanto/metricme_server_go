package main

import (
	"log"
	"time"
)

type statprocessor func(*statistic)

var gauges = make(map[string]int)
var counters = make(map[string]float32)
var timers = make(map[string][]int)
var sets = make(map[string]map[string]bool)

var flushTicker chan time.Time

var syncID int

// processCounter performs the accumulation operations on gauge metrics
func processGauge(stat *statistic) {
	switch stat.sign {
	case "+":
		gauges[stat.name] += stat.value
	case "-":
		gauges[stat.name] -= stat.value
	default:
		gauges[stat.name] = stat.value
	}
}

// processCounter performs the accumulation operations on counter metrics
func processCounter(stat *statistic) {
	counters[stat.name] += float32(stat.value) * (1 / stat.sample)
}

func processSet(stat *statistic) {
	sets[stat.name][stat.svalue] = true
}

func processTimer(stat *statistic) {
	timers[stat.name] = append(timers[stat.name], stat.value)
}

var gaugeChan = make(chan *statistic)
var timerChan = make(chan *statistic)
var counterChan = make(chan *statistic)
var setChan = make(chan *statistic)

var flushGauge = make(chan int)
var flushTimer = make(chan int)
var flushCounter = make(chan int)
var flushSet = make(chan int)

// gaugeReceiver receives parsed gauge metrics from the input channel and performs
// the required accumulation operations.
func gaugeReceiver() {
	for {
		select {
		case g := <-gaugeChan:
			processGauge(g)
		case syncVal := <-flushGauge:
			var gaugesCopy = make(map[string]int)
			for k, v := range gauges {
				gaugesCopy[k] = v
			}

			accGauge <- gaugeValues{gaugesCopy, syncVal}
		}
	}
}

// timerReceiver receives parsed timer metrics from the input channel and performs
// the required accumulation operations.
func timerReceiver() {
	for {
		select {
		case t := <-timerChan:
			processTimer(t)
		case syncVal := <-flushTimer:
			var timersCopy = make(map[string][]int)
			for k, v := range timers {
				timersCopy[k] = v
				timers[k] = make([]int, 0)
			}
			accTimer <- timerValues{timersCopy, syncVal}
		}
	}
}

// counterReceiver receives parsed counter metrics from the input channel and performs
// the required accumulation operations.
func counterReceiver() {
	for {
		select {
		case c := <-counterChan:
			processCounter(c)
		case syncVal := <-flushCounter:
			var dataCopy = make(map[string]float32)
			for k, v := range counters {
				dataCopy[k] = v
				counters[k] = 0
			}

			accCounter <- counterValues{dataCopy, 5000, syncVal}
		}
	}
}

// setReceiver receives parsed set metrics from the input channel and performs
// the required accumulation operations.
func setReceiver() {
	for {
		select {
		case s := <-setChan:
			processSet(s)
		case syncVal := <-flushSet:
			var dataCopy = make(map[string]map[string]bool)
			for k, v := range sets {
				dataCopy[k] = v
				sets[k] = make(map[string]bool)
			}
			accSet <- setValues{sets, syncVal}
		}
	}
}

// runFlushTimer starts a timer that calls the flush event every x seconds
func runFlushTimer() {
	flushTicker := time.Tick(5 * time.Second)

	for _ = range flushTicker {
		flush()
	}
}

// startFlushTimer starts a timer that calls the flush event every x seconds
func startFlushTimer() {
	go runFlushTimer()
}

// flush starts the flush operation by sending a syncID to each flush channel
func flush() {
	log.Println("Flushing")

	flushGauge <- syncID
	flushTimer <- syncID
	flushCounter <- syncID
	flushSet <- syncID

	syncID++
}

// startProcessors starts the listening loop for each different type of metric to enable
func startProcessors() {
	go gaugeReceiver()
	go timerReceiver()
	go counterReceiver()
	go setReceiver()
}

var processorChannels = map[string]chan<- *statistic{
	"g":  gaugeChan,
	"ms": timerChan,
	"c":  counterChan,
	"s":  setChan,
}
