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
	results, err := client.ReadInputRegisters(0x0000, 2)
	if err != nil {
		fmt.Println("Failed to read from SDM630 device", err)
	} else {
		fmt.Printf("L1 voltage: %.2f\n", sdm630.RtuToFloat32(results))
	}

	// TODO: Implement control channel etc, see
	// https://gist.github.com/drio/dd2c4ad72452e3c35e7e
	var rc = make(sdm630.ReadingChannel)

	qe := sdm630.NewQueryEngine(client, rc)
	td := sdm630.NewTextDumper(rc)
	go qe.Produce()
	go td.Consume()

	time.Sleep(5 * time.Second)

}
