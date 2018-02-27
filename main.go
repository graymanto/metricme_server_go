package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type statistic struct {
	name   string
	value  int
	svalue string
	mtype  string
	sign   string
	sample float32
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Metrics http handler\n")
}

func processStats(stats string) bool {
	ok, parsed := parseAllStats(stats)

	if !ok {
		return false
	}

	for _, stat := range parsed {
		processorChannels[stat.mtype] <- stat
	}

	return true
}

func acceptStat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	stats := string(body[:])
	log.Println("We received stats ", stats)

	ok := processStats(stats)

	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func main() {
	startProcessors()
	startAccumulators()
	startFlusher()
	startFlushTimer()

	router := httprouter.New()
	router.GET("/", index)
	router.POST("/stats", acceptStat)
	router.PUT("/stats", acceptStat)

	log.Fatal(http.ListenAndServe(":8080", router))
}
