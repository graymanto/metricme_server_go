package main

import (
	"log"
	"time"

	graphigo "gopkg.in/fgrosse/graphigo.v2"
)

// graphiteBackendChan is the communication channel for the console backend
var graphiteBackendChan = make(chan []statsToFlush)

func getCounterMetrics(stats statsToFlush) []graphigo.Metric {
	metrics := stats.values.(accCounterMetrics)

	toSend := make([]graphigo.Metric, 0)

	for k, val := range metrics.values {
		toSend = append(toSend, graphigo.Metric{Name: k, Value: val})
	}

	for k, val := range metrics.rates {
		toSend = append(toSend, graphigo.Metric{Name: k + ".rate", Value: val})
	}

	return toSend
}

// startGraphiteBackend starts the console backend
func startGraphiteBackend(address string, timeout int, namespace string) {
	for s := range graphiteBackendChan {
		c := graphigo.Client{
			Address: address,
			Timeout: time.Duration(timeout),
			Prefix:  namespace,
		}

		if err := c.Connect(); err != nil {
			log.Println("Unable to connect to graphite, backend not running")
			continue
		}

		defer c.Close()

		toSend := make([]graphigo.Metric, 0)

		for _, v := range s {
			switch v.statType {
			case 'c':
				cMets := getCounterMetrics(v)
				toSend = append(toSend, cMets...)
			case 't':
				log.Println("Send timers")
			case 'g':
				log.Println("Send gauges")
			case 's':
				log.Println("Send sets")
			}
		}

		c.SendAll(toSend)
	}
}
