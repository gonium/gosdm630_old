package main

import (
	"flag"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"os/signal"
)

var verbose = flag.Bool("verbose", false, "Enables extensive logging")
var brokerurl = flag.String("broker", "tcp://localhost:1883", "MQTT server url")
var username = flag.String("user", "", "Username for connecting to the MQTT server")
var password = flag.String("pass", "", "Password for connecting to the MQTT server")
var devicename = flag.String("name", "", "The name of the current measurement device")

func init() {
	flag.Parse()
	if len(*devicename) == 0 {
		log.Fatal("Please specify a name for this device (-name=<YOURID>)")
	}
}

func main() {
	for {

		// https://gist.github.com/drio/dd2c4ad72452e3c35e7e
		var rc = make(sdm630.ReadingChannel)
		var sourceControl = make(sdm630.ControlChannel)
		//var sinkControl = make(sdm630.ControlChannel)
		// handle CTRL-C correctly
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)

		// TODO: Replace q/ mqtt foo
		//qe := sdm630.NewQueryEngine(*rtuDevice, *verbose, rc, sourceControl)
		source, err := sdm630.NewMQTTSource(rc, sourceControl,
			*brokerurl, *username, *password, *devicename)
		if err != nil {
			log.Fatal("Cannot create MQTT connection: ", err)
		}
		//source := sdm630.NewTextDumper(rc, sinkControl)
		//source := sdm630.NewTextGui(rc, sinkControl)
		//go qe.Produce()
		go source.Run()
		select {
		case _ = <-signals:
			log.Fatal("received SIGTERM, exiting.")
			break
			//case pm := <-sourceControl:
			//	if pm == sdm630.ControlReadFailure {
			//		// TODO: Collect statistics here.
			//		log.Println("Read Failure")
			//	} else if pm == sdm630.ControlClose {
			//		log.Println("Source closed.")
			//	} else {
			//		log.Fatal("Unknown control message from source:", pm)
			//	}
			//case tm := <-sinkControl:
			//	log.Fatal(tm)
			//	break
		}
	}
}
