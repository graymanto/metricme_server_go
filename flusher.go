package main

import "log"

var flushValues = make(map[int][]statsToFlush)

type backendChannel chan []statsToFlush

func flushProcessor(backends []backendChannel) {
	for v := range flusher {
		flushValues[v.syncID] = append(flushValues[v.syncID], v)

		if len(flushValues[v.syncID]) == 4 {
			for _, bechan := range backends {
				bechan <- flushValues[v.syncID]
			}

			delete(flushValues, v.syncID)
			for k := range flushValues {
				log.Println("Remaining key is", k)
			}
		}
	}
}

func startFlusher(config *systemConfig) {
	backends := make([]backendChannel, 0)

	if config.flushConsole {
		go startConsoleBackend()
		backends = append(backends, consoleBackendChan)
	}

	if config.flushGraphite {
		go startGraphiteBackend(config.graphiteAddress, config.graphiteTimeout,
			config.graphiteNamespace)
		backends = append(backends, graphiteBackendChan)
	}

	go flushProcessor([]backendChannel{consoleBackendChan})
}
