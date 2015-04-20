package sdm630

import (
	"encoding/binary"
	"math"
)

func RtuToFloat32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
