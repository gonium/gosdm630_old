package main

import (
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"time"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")

func init() {
	flag.Parse()
}

func main() {
	for {
		// Modbus RTU/ASCII
		handler := modbus.NewRTUClientHandler(*rtuDevice)
		handler.BaudRate = 9600
		handler.DataBits = 8
		handler.Parity = "N"
		handler.StopBits = 1
		handler.SlaveId = 1
		handler.Timeout = 5 * time.Second
		if *verbose {
			handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
			fmt.Printf("Connecting to RTU via %s\n", *rtuDevice)
		}

		err := handler.Connect()
		if err != nil {
			log.Fatal("Failed to connect: ", err)
		}
		defer handler.Close()

		client := modbus.NewClient(handler)

		// https://gist.github.com/drio/dd2c4ad72452e3c35e7e
		var rc = make(sdm630.ReadingChannel)
		var producerControl = make(sdm630.ControlChannel)
		var consumerControl = make(sdm630.ControlChannel)

		qe := sdm630.NewQueryEngine(client, rc, producerControl)
		//td := sdm630.NewTextDumper(rc, consumerControl)
		td := sdm630.NewTextGui(rc, consumerControl)
		go qe.Produce()
		go td.ConsumeData()
		// TODO: Select over control channels, restart serial interface in
		// case of failures.
		select {
		case pm := <-producerControl:
			if pm == sdm630.ControlReadFailure {
				log.Printf("Read Failure")
			}
			if pm == sdm630.ControlClose {
				log.Printf("Producer closed.")
			} else {
				log.Fatal("Unknown control message from producer:", pm)
			}
		case tm := <-consumerControl:
			log.Fatal(tm)
			break
		}
	}
}
