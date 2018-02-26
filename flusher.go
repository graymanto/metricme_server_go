package main

import "log"

var flushValues = make(map[int][]statsToFlush)

func flushProcessor() {
	for v := range flusher {
		log.Println("Ready to flush", v)

		flushValues[v.syncID] = append(flushValues[v.syncID], v)

		if len(flushValues[v.syncID]) == 4 {
			log.Println("We have everything, we should flush!!")
			delete(flushValues, v.syncID)
		}
	}
}

func startFlusher() {
	go flushProcessor()
}
