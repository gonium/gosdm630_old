package sdm630

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"math"
	"time"
)

const (
	OpCodeL1Voltage     = 0x0000
	OpCodeL2Voltage     = 0x0002
	OpCodeL3Voltage     = 0x0004
	OpCodeL1Current     = 0x0006
	OpCodeL2Current     = 0x0008
	OpCodeL3Current     = 0x000A
	OpCodeL1Power       = 0x000C
	OpCodeL2Power       = 0x000E
	OpCodeL3Power       = 0x0010
	OpCodeL1PowerFactor = 0x001e
	OpCodeL2PowerFactor = 0x0020
	OpCodeL3PowerFactor = 0x0022
)

type QueryEngine struct {
	client     modbus.Client
	datastream ReadingChannel
	control    ControlChannel
}

func NewQueryEngine(handler modbus.Client, channel ReadingChannel,
	c ControlChannel) *QueryEngine {
	return &QueryEngine{client: handler, datastream: channel, control: c}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float32,
	err error) {
	results, err := q.client.ReadInputRegisters(opcode, 2)
	if err == nil {
		retval = RtuToFloat32(results)
	}
	return retval, err
}

func (q *QueryEngine) queryOrFail(opcode uint16) (retval float32) {
	retval, err := q.retrieveOpCode(opcode)
	if err != nil {
		q.control <- ControlReadFailure
	}
	return
}

func (q *QueryEngine) Produce() {
	// First: Query the SDM630 device for all interesting data.
	for {
		q.datastream <- Readings{
			L1Voltage: q.queryOrFail(OpCodeL1Voltage),
			L2Voltage: q.queryOrFail(OpCodeL2Voltage),
			L3Voltage: q.queryOrFail(OpCodeL3Voltage),
			L1Current: q.queryOrFail(OpCodeL1Current),
			L2Current: q.queryOrFail(OpCodeL2Current),
			L3Current: q.queryOrFail(OpCodeL3Current),
			L1Power:   q.queryOrFail(OpCodeL1Power),
			L2Power:   q.queryOrFail(OpCodeL2Power),
			L3Power:   q.queryOrFail(OpCodeL3Power),
			L1CosPhi:  q.queryOrFail(OpCodeL1PowerFactor),
			L2CosPhi:  q.queryOrFail(OpCodeL2PowerFactor),
			L3CosPhi:  q.queryOrFail(OpCodeL3PowerFactor),
		}
		time.Sleep(1 * time.Second)
	}
	q.control <- ControlClose
}

func RtuToFloat32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
