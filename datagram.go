package sdm630

import (
	"fmt"
)

const (
	ControlClose       = iota
	ControlReadFailure = iota
)

type ControlChannel chan int64

type ReadingChannel chan Readings

type Readings struct {
	L1Voltage float32
	L2Voltage float32
	L3Voltage float32
	L1Current float32
	L2Current float32
	L3Current float32
	L1Power   float32
	L2Power   float32
	L3Power   float32
	L1CosPhi  float32
	L2CosPhi  float32
	L3CosPhi  float32
}

func (r *Readings) String() string {
	fmtString := "L1: %.2fV %.2fA %.2fW %.2fcos | " +
		"L2: %.2fV %.2fA %.2fW %.2fcos | " +
		"L3: %.2fV %.2fA %.2fW %.2fcos"
	return fmt.Sprintf(fmtString,
		r.L1Voltage,
		r.L1Current,
		r.L1Power,
		r.L1CosPhi,
		r.L2Voltage,
		r.L2Current,
		r.L2Power,
		r.L2CosPhi,
		r.L3Voltage,
		r.L3Current,
		r.L3Power,
		r.L3CosPhi,
	)
}
