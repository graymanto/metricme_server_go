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

func makeParseTestString(name string, value int, mtype string, sign string,
	sample string) string {
	var tvalue string = sign + strconv.Itoa(value)
	testString := name + ":" + tvalue + "|" + mtype
	if sample != "" {
		testString += "|" + sample
	}

	return testString
}

func makeParseTestCase(name string, value int, mtype string, sign string,
	sample string) *parseStatTest {

	testString := makeParseTestString(name, value, mtype, sign, sample)

	test := &parseStatTest{testString, []*statistic{}}
	validSample := sample
	if mtype != "c" && mtype != "ms" {
		validSample = ""
	}
	test.out = append(test.out, &statistic{name, value, mtype, sign, validSample})

	return test
}

func appendParseTestCase(test *parseStatTest, name string, value int, mtype string,
	sign string, sample string) {
	testString := makeParseTestString(name, value, mtype, sign, sample)
	test.in += ";" + testString
	test.out = append(test.out, &statistic{name, value, mtype, sign, sample})
}

func newParseStatTests() []*parseStatTest {
	var tests = []*parseStatTest{}

	tests = append(tests, makeParseTestCase("gagey", 333, "g", "", ""))
	tests = append(tests, makeParseTestCase("counter", 17, "c", "", ""))

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

	return tests
}

func TestParseStat(t *testing.T) {
	var tests = newParseStatTests()

	for _, tt := range tests {
		ok, stats := parseAllStats(tt.in)

		if tt.out == nil {
			if stats != nil || ok {
				t.Error("Stat parsing has not failed as expected", tt.in, ok, stats == nil)
			}
		} else {
			if !ok || stats == nil {
				t.Error("Stat not parsed as expected", tt.in, ok, stats == nil)
				continue
			}

			for i, stat := range stats {
				if !reflect.DeepEqual(stat, tt.out[i]) {
					t.Error("Parse result not as expected", stat, tt.out[i])
				}
			}
		}
	}
}
