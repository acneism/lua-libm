package libm

import (
	"encoding/binary"
	"math"
)

func tableValue(bits string, i int) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64([]byte(bits[i*8 : i*8+8])))
}

func highWord(x float64) int {
	return int(math.Float64bits(x) >> 32)
}

func lowWord(x float64) int {
	return int(uint32(math.Float64bits(x)))
}

func exactAdd(x, y float64) (z, zz float64) {
	z = float64(x + y)
	if math.Abs(x) > math.Abs(y) {
		zz = float64(float64(x-z) + y)
	} else {
		zz = float64(float64(y-z) + x)
	}
	return
}

func exactSub(x, y float64) (z, zz float64) {
	z = float64(x - y)
	if math.Abs(x) > math.Abs(y) {
		zz = float64(float64(x-z) - y)
	} else {
		zz = float64(x - float64(y+z))
	}
	return
}

func dividePair(x, xx, y, yy float64) (z, zz float64) {
	c := float64(x / y)
	u := float64(c * y)
	uu := math.FMA(c, y, -u)
	cc := math.FMA(-c, yy, float64(float64(float64(x-u)-uu)+xx)) / y
	z = float64(c + cc)
	zz = float64(float64(c-z) + cc)
	return
}
