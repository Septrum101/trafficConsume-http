package main

import (
	"fmt"
)

var (
	appName = "TrafficConsume-HTTP"
	version = "dev"
	date    = "unknown"
)

func getVersion() {
	fmt.Printf("%s %s, built at %s\n", appName, version, date)
}
