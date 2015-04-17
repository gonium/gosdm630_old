package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gonium/modbus"
	"log"
	"math"
	"os"
	"time"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")

func init() {
	flag.Parse()
}

func main() {
	fmt.Printf("Connecting to RTU via %s\n", *rtuDevice)
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler(*rtuDevice)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)

	err := handler.Connect()
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadInputRegisters(0, 2)
	if err != nil {
		fmt.Println("Failed to read from SDM630 device", err)
	} else {
		l1Bits := binary.BigEndian.Uint32(results)
		l1Float := math.Float32frombits(l1Bits)
		fmt.Printf("L1 voltage: %.2f\n", l1Float)
	}
}
