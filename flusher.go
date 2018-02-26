package main

import "log"

var flushValues = make(map[int][]statsToFlush)

func flushProcessor() {
	for v := range flusher {
		log.Println("Received flush with id", v.syncID, string(v.statType))

		flushValues[v.syncID] = append(flushValues[v.syncID], v)

		if len(flushValues[v.syncID]) == 4 {
			delete(flushValues, v.syncID)
			log.Println("We have everything, we should flush!!", v.syncID, len(flushValues))
			log.Println("Flushing", flushValues[v.syncID])
			for k := range flushValues {
				log.Println("Remaining key is", k)
			}
		}
	}
}

func startFlusher() {
	go flushProcessor()
}
