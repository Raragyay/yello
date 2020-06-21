package main

import (
	"fmt"
	"strings"
)

type parseFlag int

const (
	ok                 parseFlag = 0
	invalidParamsCount parseFlag = 1
	notPONG            parseFlag = 2
)

//parseUtilsAndSignal parses the given string in line with the MD protocol. if expectedLen is 0 it is disregarded and any number of fields more than 0 are accepted.
func parseUtilsAndSignal(text string, expectedLen int) ([]string, parseFlag) {
	fields := strings.Fields(text)
	flen := len(fields)
	if flen == 0 {
		return nil, invalidParamsCount
	}
	if expectedLen != 0 {
		if flen != expectedLen {
			return nil, invalidParamsCount
		}
	}
	if fields[0] != "PONG" {
		return nil, notPONG
	}

	return fields, ok
}

func handlepanic() {

	if a := recover(); a != nil {

		fmt.Println("RECOVER", a)
	}
}

func parseFlagToString(pflag parseFlag) string {
	switch pflag {
	case ok:
		return "ok"
	case notPONG:
		return "notPONG"
	case invalidParamsCount:
		return "invalidParamsCount"
	}
	return ""
}
