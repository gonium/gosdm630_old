package main

import (
	"flag"
	"fmt"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"os/signal"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")
var broker = flag.String("broker", "localhost:1883", "MQTT server address")
var username = flag.String("user", "", "Username for connecting to the MQTT server")
var password = flag.String("pass", "", "Password for connecting to the MQTT server")
var devicename = flag.String("name", "", "The name of the current measurement device")
var interval = flag.Int("interval", 5, "Seconds between querying the SDM630")

func init() {
	flag.Parse()
	if len(*devicename) == 0 {
		log.Fatal("Please specify a name for this device (-name=<YOURID>)")
	}
}

func main() {
	var rc = make(sdm630.ReadingChannel)
	var producerControl = make(sdm630.ControlChannel)
	var consumerControl = make(sdm630.ControlChannel)
	// handle CTRL-C correctly
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	qe := sdm630.NewQueryEngine(*rtuDevice, *interval, *verbose, rc, producerControl)
	topic := fmt.Sprintf("readings/%s", *devicename)
	td, err := sdm630.NewMQTTSubmitter(rc, consumerControl,
		*broker, *username, *password, topic)
	if err != nil {
		log.Fatal("Cannot create MQTT connection: ", err)
	}
	defer td.Close()
	go td.ConsumeData()
	go qe.Produce()
	//td := sdm630.NewTextDumper(rc, consumerControl)
	//td := sdm630.NewTextGui(rc, consumerControl)
	for {
		// TODO: Select over control channels, restart serial interface in
		// case of failures.
		select {
		case _ = <-signals:
			log.Fatal("received SIGTERM, exiting.")
			break
		case pm := <-producerControl:
			if pm == sdm630.ControlReadFailure {
				// TODO: Collect statistics here.
				log.Println("SDM630: Cannot read from device (no response?)")
			} else if pm == sdm630.ControlClose {
				log.Println("Producer closed, exiting.")
				break
			} else {
				log.Fatal("Unknown control message from producer:", pm)
			}
		case tm := <-consumerControl:
			log.Fatal(tm)
			break
		}
	}
}
