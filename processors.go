package main

import (
	"log"
	"time"
)

type statprocessor func(*statistic)

var gauges = make(map[string]int)
var counters = make(map[string]int)
var timers = make(map[string]int)
var sets = make(map[string]int)

var flushTicker chan time.Time

var syncID int

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

func processCounter(stat *statistic) {
	// TODO: sampling
	counters[stat.name] += stat.value
}

var gaugeChan = make(chan *statistic)
var timerChan = make(chan *statistic)
var counterChan = make(chan *statistic)
var setChan = make(chan *statistic)

var flushGauge = make(chan int)
var flushTimer = make(chan int)
var flushCounter = make(chan int)
var flushSet = make(chan int)

func gaugeReceiver() {
	for {
		select {
		case g := <-gaugeChan:
			processGauge(g)
			log.Println("Processed gauges", gauges)
		case syncVal := <-flushGauge:
			var gaugesCopy = make(map[string]int)
			for k, v := range gauges {
				gaugesCopy[k] = v
			}

			accGauge <- gaugeValues{gaugesCopy, syncVal}
		}
	}
}

func timerReceiver() {
	for {
		select {
		case t := <-timerChan:
			log.Println("Received timers", t)
		case syncVal := <-flushTimer:
			accTimer <- timerValues{timers, syncVal}
		}
	}
}

func counterReceiver() {
	for {
		select {
		case c := <-counterChan:
			log.Println("Received counter", c)
			processCounter(c)
		case syncVal := <-flushCounter:
			var dataCopy = make(map[string]int)
			for k, v := range counters {
				dataCopy[k] = v
				counters[k] = 0
			}

			accCounter <- counterValues{dataCopy, syncVal}
		}
	}
}

func setReceiver() {
	for {
		select {
		case s := <-setChan:
			log.Println("Received set", s)
		case syncVal := <-flushSet:
			accSet <- setValues{sets, syncVal}
		}
	}
}

func processFlushTimer() {
	flushTicker := time.Tick(5 * time.Second)

	for _ = range flushTicker {
		flush()
	}
}

func startFlushTimer() {
	go processFlushTimer()
}

func flush() {
	log.Println("Flushing id", syncID)

	flushGauge <- syncID
	flushTimer <- syncID
	flushCounter <- syncID
	flushSet <- syncID

	syncID++
}

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
