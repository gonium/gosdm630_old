package sdm630

import (
	"fmt"
)

type Readings struct {
	Error     error
	L1Voltage float32
	L2Voltage float32
	L3Voltage float32
}

func (r *Readings) String() string {
	return fmt.Sprintf("Voltages: L1=%.2f L2=%.2f L3=%.2f")
}
