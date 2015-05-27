package main

import (
	"flag"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"os/signal"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")

func init() {
	flag.Parse()
}

func main() {
	for {

		// https://gist.github.com/drio/dd2c4ad72452e3c35e7e
		var rc = make(sdm630.ReadingChannel)
		var producerControl = make(sdm630.ControlChannel)
		var consumerControl = make(sdm630.ControlChannel)
		// handle CTRL-C correctly
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)

		qe := sdm630.NewQueryEngine(*rtuDevice, *verbose, rc, producerControl)
		td := sdm630.NewTextDumper(rc, consumerControl)
		//td := sdm630.NewTextGui(rc, consumerControl)
		go qe.Produce()
		go td.ConsumeData()
		// TODO: Select over control channels, restart serial interface in
		// case of failures.
		select {
		case _ = <-signals:
			log.Fatal("received SIGTERM, exiting.")
			break
		case pm := <-producerControl:
			if pm == sdm630.ControlReadFailure {
				log.Println("Read Failure")
			} else if pm == sdm630.ControlClose {
				log.Println("Producer closed.")
			} else {
				log.Fatal("Unknown control message from producer:", pm)
			}
		case tm := <-consumerControl:
			log.Fatal(tm)
			break
		}
	}
}
