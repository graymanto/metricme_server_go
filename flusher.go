package main

import "log"

var flushValues = make(map[int][]statsToFlush)

func flushProcessor() {
	for v := range flusher {
		flushValues[v.syncID] = append(flushValues[v.syncID], v)

		if len(flushValues[v.syncID]) == 4 {
			consoleBackendChan <- flushValues[v.syncID]

			delete(flushValues, v.syncID)
			for k := range flushValues {
				log.Println("Remaining key is", k)
			}
		}
	}
}

func startFlusher() {
	go flushProcessor()
	go startConsoleBackend()
}
