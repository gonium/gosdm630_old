package sdm630

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ReadingChannel chan Readings

type Readings struct {
	Time      time.Time
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
	fmtString := "T: %s - L1: %.2fV %.2fA %.2fW %.2fcos | " +
		"L2: %.2fV %.2fA %.2fW %.2fcos | " +
		"L3: %.2fV %.2fA %.2fW %.2fcos"
	return fmt.Sprintf(fmtString,
		r.Time.Format(time.RFC3339),
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

func (r *Readings) JSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
