package main

import (
	"log"
	"strconv"
	"strings"
)

type statparser func(string, []string) (bool, *statistic)

func newBasicStat(name string, parts []string) (bool, *statistic) {
	num, err := strconv.Atoi(parts[0])

	if err != nil {
		log.Println("Error parsing number from stat", name, parts[0])
		return false, nil
	}

	return true, &statistic{
		name, num, parts[1], "", "",
	}
}

func parseCounting(name string, parts []string) (bool, *statistic) {
	ok, basic := newBasicStat(name, parts)

	if !ok {
		return false, nil
	}

	if len(parts) < 3 || !strings.HasPrefix(parts[2], "@") {
		return true, basic
	}

	basic.sample = parts[2]

	return true, basic
}

func parseTiming(name string, parts []string) (bool, *statistic) {
	return newBasicStat(name, parts)
}

func parseGauges(name string, parts []string) (bool, *statistic) {
	return newBasicStat(name, parts)
}

var parseMap = map[string]statparser{
	"c":  parseCounting,
	"ms": parseTiming,
	"g":  parseTiming,
}

func parseStat(stat string) (bool, *statistic) {
	parts := strings.Split(stat, ":")

	if len(parts) != 2 {
		log.Println("ERROR invalid stat parsed", stat)
		return false, nil
	}

	body := strings.Split(parts[1], "|")
	if len(body) < 2 {
		log.Println("ERROR invalid stat body", stat)
		return false, nil
	}

	if parser, ok := parseMap[body[1]]; ok {
		return parser(parts[0], body)
	}

	log.Println("ERROR unknown stat type parsed", stat)
	return false, nil
}

func parseAllStats(stats string) (bool, []*statistic) {
	splitStats := strings.Split(stats, ";")

	parsed := make([]*statistic, len(splitStats))

	for i, stat := range splitStats {
		ok, stat := parseStat(stat)

		if !ok {
			return false, nil
		}

		parsed[i] = stat
	}

	return true, parsed
}
