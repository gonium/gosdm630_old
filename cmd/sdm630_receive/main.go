package main

import (
	"flag"
	"fmt"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"os/signal"
	"strings"
)

var verbose = flag.Bool("verbose", false, "Enables extensive logging")
var brokerurl = flag.String("broker", "localhost:1883", "MQTT server address")
var username = flag.String("user", "", "Username for connecting to the MQTT server")
var password = flag.String("pass", "", "Password for connecting to the MQTT server")
var devicename = flag.String("name", "", "The name of the current measurement device")

func init() {
	flag.Parse()
	if len(*devicename) == 0 {
		log.Fatal("Please specify a name for this device (-name=<YOURID>)")
	}
}

func printReadings(topicName []byte, message []byte) {
	s := strings.Split(string(topicName), "/")
	msgtype, devicename, measurement, subcategory := s[0], s[1], s[2],
		s[3]
	switch msgtype {
	case "readings":
		log.Printf("%s: %s(%s) = %s", devicename, measurement, subcategory,
			string(message))
		break
	default:
		log.Println("unknown message type, topic was ", string(topicName))
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

		source, err := sdm630.NewMQTTSource(rc, sourceControl,
			*brokerurl, *username, *password, *devicename)
		defer source.Close()
		if err != nil {
			log.Fatal("Cannot create MQTT connection: ", err)
		}
		topic := fmt.Sprintf("readings/%s", *devicename)
		source.Subscribe(topic, printReadings)
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
