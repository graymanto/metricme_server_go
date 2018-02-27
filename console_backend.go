package main

import "log"

// consoleBackendChan is the communication channel for the console backend
var consoleBackendChan = make(chan []statsToFlush)

func printCounters(stats statsToFlush) {
	metrics := stats.values.(accCounterMetrics)

	for k, val := range metrics.values {
		log.Println("Counter:", k, val)
	}

	for k, val := range metrics.rates {
		log.Println("Counter rate:", k, val)
	}
}

func printTimers(stats statsToFlush) {
	// Not implemented printing timers
}

func printGauges(stats statsToFlush) {
	metrics := stats.values.(map[string]int)

	for k, val := range metrics {
		log.Println("Gauge:", k, val)
	}
}

func printSets(stats statsToFlush) {
	metrics := stats.values.(map[string]map[string]bool)

	for k, val := range metrics {
		log.Println("Set:", k, val)
	}
}

// startConsoleBackend starts the console backend
func startConsoleBackend() {
	for s := range consoleBackendChan {
		// TODO: sort slice first for consistent console printing
		for _, v := range s {
			switch v.statType {
			case 'c':
				printCounters(v)
			case 't':
				printTimers(v)
			case 'g':
				printGauges(v)
			case 's':
				printSets(v)
			}
		}
	}
}
