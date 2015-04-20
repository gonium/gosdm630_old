package sdm630

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"math"
)

const (
	OpCodeL1Voltage = 0x0000
	OpCodeL2Voltage = 0x0002
	OpCodeL3Voltage = 0x0004
)

type ReadingChannel chan Readings

type QueryEngine struct {
	client     modbus.Client
	datastream ReadingChannel
}

func NewQueryEngine(handler modbus.Client, channel ReadingChannel) *QueryEngine {
	return &QueryEngine{client: handler, datastream: channel}
}

func (q *QueryEngine) Produce() {
	for i := 0; i < 10; i++ {
		q.datastream <- Readings{
			Error:     nil,
			L1Voltage: 230.1,
			L2Voltage: 230.2,
			L3Voltage: 230.3,
		}
	}
}

func RtuToFloat32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
