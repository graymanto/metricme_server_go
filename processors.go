package main

import "log"

type statprocessor func(*statistic)

var gauges = make(map[string]int)
var counters = make(map[string]int)

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
		case syncID := <-flushGauge:
			log.Println("We should be flushing gauge!!")
			var gaugesCopy = make(map[string]int)
			for k, v := range gauges {
				gaugesCopy[k] = v
			}

			accGauge <- gaugeValues{gaugesCopy, syncID}
		}
	}
}

func timerReceiver() {
	for {
		select {
		case t := <-timerChan:
			log.Println("Received timer", t)
		case <-flushTimer:
			log.Println("We should be flushing timer!!")
		}
	}
}

func counterReceiver() {
	for {
		select {
		case c := <-counterChan:
			log.Println("Received counter", c)
			processCounter(c)
		case syncID := <-flushCounter:
			log.Println("We should be flushing counter!!")

			var dataCopy = make(map[string]int)
			for k, v := range counters {
				dataCopy[k] = v
				counters[k] = 0
			}

			accCounter <- counterValues{dataCopy, syncID}
		}
	}
}

func setReceiver() {
	for {
		select {
		case s := <-setChan:
			log.Println("Received set", s)
		case <-flushSet:
			log.Println("We should be flushing sets!!")
		}
	}
}

func flush() {
	flushGauge <- syncID
	flushTimer <- syncID
	flushCounter <- syncID
	flushSet <- syncID
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
