package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/Septrum101/trafficConsume-http/app"
)

func main() {
	getVersion()

	//logger := log.StandardLogger()
	//logger.SetLevel(log.DebugLevel)

	cli := app.New()
	if err := cli.Start(); err != nil {
		log.Panic(err)
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-osSignals
	cli.Close()
}
