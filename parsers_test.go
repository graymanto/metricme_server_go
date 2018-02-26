package main

import (
	"reflect"
	"strconv"
	"testing"
)

type parseStatTest struct {
	in  string
	out []*statistic
}

// Makes a statsd compatible metric string from the passed in values.
func makeParseTestString(name string, value int, mtype string, sign string,
	sample string) string {
	var tvalue string = sign + strconv.Itoa(value)
	testString := name + ":" + tvalue + "|" + mtype
	if sample != "" {
		testString += "|" + sample
	}

	return testString
}

// Creates a single test case. Builds a statsd compatible string and also creates
// the expected test result in a statistic struct.
func makeParseTestCase(name string, value int, mtype string, sign string,
	sample string) *parseStatTest {

	testString := makeParseTestString(name, value, mtype, sign, sample)

	test := &parseStatTest{testString, []*statistic{}}
	validSample := sample

	// All types except counters and timers should ignore sampling
	if mtype != "c" && mtype != "ms" {
		validSample = ""
	}
	test.out = append(test.out, &statistic{name, value, mtype, sign, validSample})

	return test
}

// Appends a new test case to an existing one. Appends to the statsd compatible string
// and and appends a new result struct to the test case.
func appendParseTestCase(test *parseStatTest, name string, value int, mtype string,
	sign string, sample string) {
	testString := makeParseTestString(name, value, mtype, sign, sample)
	test.in += ";" + testString
	test.out = append(test.out, &statistic{name, value, mtype, sign, sample})
}

// Builds a table of test cases for the parse all stats tests
func newParseStatsTests() []*parseStatTest {
	var tests = []*parseStatTest{}

	tests = append(tests, makeParseTestCase("gagey", 333, "g", "", ""))
	tests = append(tests, makeParseTestCase("counter", 17, "c", "", ""))
	tests = append(tests, makeParseTestCase("counter", 19, "c", "", "@0.5"))
	tests = append(tests, makeParseTestCase("gagey", 35, "g", "", "@0.5"))
	tests = append(tests, makeParseTestCase("timer", 17, "ms", "", ""))
	tests = append(tests, makeParseTestCase("gaugeSigned", 35, "g", "+", ""))
	tests = append(tests, makeParseTestCase("gaugeSigned", 35, "g", "-", ""))

	multi := makeParseTestCase("gauge1", 11, "g", "", "")
	appendParseTestCase(multi, "counter1", 12, "c", "", "@0.1")

	tests = append(tests, multi)

	multi = makeParseTestCase("gauge", 33, "g", "", "")
	appendParseTestCase(multi, "counter", 55, "c", "", "")
	appendParseTestCase(multi, "gauge", 22, "g", "", "")

	tests = append(tests, multi)

	tests = append(tests, &parseStatTest{"", nil})
	tests = append(tests, &parseStatTest{"invalid1", nil})
	tests = append(tests, &parseStatTest{"invalid2:", nil})
	tests = append(tests, &parseStatTest{"invalid3:5", nil})
	tests = append(tests, &parseStatTest{"rubbishtype:5|q", nil})
	tests = append(tests, &parseStatTest{"badNumber:a|c", nil})
	tests = append(tests, &parseStatTest{"badGauge:x35|g", nil})

	// One bad stat in multi test case, should fail everything
	tests = append(tests, &parseStatTest{"gauge:1|g;count:1|c;bad:6", nil})

	return tests
}

func TestParseAllStats(t *testing.T) {
	var tests = newParseStatsTests()

	for _, tt := range tests {
		ok, stats := parseAllStats(tt.in)

		// Check failure test cases were processed correctly
		if tt.out == nil {
			if stats != nil || ok {
				t.Error("Stat parsing has not failed as expected", tt.in, ok, stats == nil)
			}
		} else {
			// Check correct test cases didn't fail parsing
			if !ok || stats == nil {
				t.Error("Stat not parsed as expected", tt.in, ok, stats == nil)
				continue
			}

			// Check correct test cases parsed as expected
			for i, stat := range stats {
				if !reflect.DeepEqual(stat, tt.out[i]) {
					t.Error("Parse result not as expected", stat, tt.out[i])
				}
			}
		}
	}
}
