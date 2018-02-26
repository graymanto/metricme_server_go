package main

type statprocessor func(*statistic)

var gauges = make(map[string]int)

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

var gaugeChan = make(chan *statistic)

func gaugeReceiver() {
	for gauge := range gaugeChan {
		processGauge(gauge)
	}
}

var processorChannels = map[string]chan<- *statistic{
	"g": gaugeChan,
}
