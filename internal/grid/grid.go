package grid

import (
	"math"
	"strconv"

	libm "github.com/acneism/lua-libm"
)

const BaseCount = 100000
const boundaryMultipleLimit = 4096
const neighborsPerCenter = 5
const boundaryMultiples = 2*boundaryMultipleLimit + 1
const functionSeedStep = uint64(0x100000001b3)

type Function struct {
	Name     string
	Evaluate func(float64, float64) float64
	Generate func(int) (float64, float64)
	Count    int
}

func unary(operation func(float64) float64) func(float64, float64) float64 {
	return func(x, _ float64) float64 { return operation(x) }
}

func Functions() []Function {
	periods := []float64{math.Pi / 2, math.Pi}
	tangentPeriods := []float64{math.Pi / 2, math.Pi, math.Pi / 4}
	return []Function{
		{"sin", unary(libm.Sin), circularGenerator(0, periods), BaseCount + len(periods)*boundaryMultiples*neighborsPerCenter},
		{"cos", unary(libm.Cos), circularGenerator(functionSeedStep, periods), BaseCount + len(periods)*boundaryMultiples*neighborsPerCenter},
		{"tan", unary(libm.Tan), circularGenerator(2*functionSeedStep, tangentPeriods), BaseCount + len(tangentPeriods)*boundaryMultiples*neighborsPerCenter},
		{"asin", unary(libm.Asin), inverseGenerator(3 * functionSeedStep), BaseCount + 2*neighborsPerCenter},
		{"acos", unary(libm.Acos), inverseGenerator(4 * functionSeedStep), BaseCount + 2*neighborsPerCenter},
		{"atan", unary(libm.Atan), ordinaryGenerator(5 * functionSeedStep), BaseCount},
		{"atan2", libm.Atan2, atan2Input, BaseCount},
	}
}

func mix(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

func unit(x uint64) float64 { return float64(float64(x>>11) * 0x1p-53) }

func signed(x float64, bits uint64) float64 {
	if bits&1 != 0 {
		return -x
	}
	return x
}

func special(i uint64) float64 {
	bits := [12]uint64{0, 1 << 63, 1, 1<<63 | 1, 0x000fffffffffffff, 0x0010000000000000, 0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000000, 0xfff8000000000000, 0x3ff0000000000000, 0xbff0000000000000}
	return math.Float64frombits(bits[i%uint64(len(bits))])
}

func seeds(seed uint64, index int) (uint64, uint64, uint64) {
	a := mix(uint64(index) + seed)
	b := mix(a)
	return a, b, mix(b)
}

func ordinaryGenerator(seed uint64) func(int) (float64, float64) {
	return func(index int) (float64, float64) { return ordinaryInput(seed, index) }
}

func circularGenerator(seed uint64, periods []float64) func(int) (float64, float64) {
	return func(index int) (float64, float64) {
		if index < BaseCount {
			return ordinaryInput(seed, index)
		}
		boundary := index - BaseCount
		period := periods[boundary/(boundaryMultiples*neighborsPerCenter)]
		multiple := boundary/neighborsPerCenter%boundaryMultiples - boundaryMultipleLimit
		center := float64(float64(multiple) * period)
		return neighbor(center, boundary%neighborsPerCenter-2), 0
	}
}

func inverseGenerator(seed uint64) func(int) (float64, float64) {
	return func(index int) (float64, float64) {
		if index < BaseCount {
			return inverseInput(seed, index)
		}
		boundary := index - BaseCount
		center := 1.0
		if boundary/neighborsPerCenter == 0 {
			center = -1
		}
		return neighbor(center, boundary%neighborsPerCenter-2), 0
	}
}

func neighbor(center float64, offset int) float64 {
	direction := math.Inf(1)
	if offset < 0 {
		direction = math.Inf(-1)
		offset = -offset
	}
	for range offset {
		center = math.Nextafter(center, direction)
	}
	return center
}

func ordinaryInput(seed uint64, index int) (x, y float64) {
	if index < 144 {
		return special(uint64(index)), special(uint64(index / 12))
	}
	a, b, c := seeds(seed, index)
	y = math.Float64frombits(b)
	switch index % 8 {
	case 0:
		x = float64(unit(a)*64) - 32
	case 1:
		x = float64(unit(a)*2e8) - 1e8
	case 2:
		x = math.Float64frombits(a)
	case 3:
		x = math.Ldexp(signed(1+unit(a), b), int(c%2098)-1074)
	case 4:
		x = float64(float64(int64(a%200000001)-100000000) * (math.Pi / 2))
		x = math.Nextafter(x, signed(math.Inf(1), b))
	case 5:
		x = signed(float64(a%256)/256+math.Ldexp(unit(b), -16), c)
	case 6:
		x = [12]float64{0x1p-27, 0x1p-26, 0.126, 0.85546875, 2.426265, 105414350, 0.0608, 0.787, 25, 1e8, 0.0625, 16}[a%12]
		x = signed(math.Nextafter(x, signed(math.Inf(1), b)), c)
	case 7:
		x = signed(math.Ldexp(unit(a), -int(b%1074)), c)
	}
	return x, y
}

func inverseInput(seed uint64, index int) (x, y float64) {
	if index < 144 {
		return special(uint64(index)), special(uint64(index / 12))
	}
	a, b, c := seeds(seed, index)
	y = math.Float64frombits(b)
	switch index % 8 {
	case 0, 7:
		x = float64(unit(a)*2) - 1
	case 1:
		x = signed(1-math.Ldexp(unit(a), -int(b%52)), c)
	case 2:
		x = math.Float64frombits(a)
	case 3:
		x = signed(float64(a%128)/128+math.Ldexp(unit(b), -12), c)
	case 4:
		x = [8]float64{0.125, 0.25, 0.5, 0.75, 0.921875, 0.953125, 0.96875, 1}[a%8]
		x = signed(math.Nextafter(x, signed(math.Inf(1), b)), c)
	case 5:
		x = math.Float64frombits(a & 0x800fffffffffffff)
	case 6:
		x = signed(math.Nextafter(1, 2), a)
	}
	return x, y
}

func atan2Input(index int) (float64, float64) {
	if index < 144 {
		return special(uint64(index)), special(uint64(index / 12))
	}
	a, b, c := seeds(6*functionSeedStep, index)
	switch index % 8 {
	case 0:
		return float64(unit(a)*64) - 32, float64(unit(b)*64) - 32
	case 1:
		return signed(unit(a), b), signed(1, c)
	case 2:
		return math.Float64frombits(a), math.Float64frombits(b)
	case 3:
		exponent := int(c%2100) - 1074
		return math.Ldexp(signed(unit(a), b), exponent), math.Ldexp(signed(unit(b), a), exponent)
	case 4:
		return special(a), special(b)
	case 5:
		return signed(1+float64(unit(a)*15), b), signed(1+float64(unit(b)*15), c)
	case 6:
		return math.Ldexp(signed(1+unit(a), b), int(c%120)-60), signed(1+unit(b), a)
	default:
		return math.Float64frombits(a & 0x800fffffffffffff), math.Float64frombits(b & 0x800fffffffffffff)
	}
}

func Argument(x float64) string {
	if math.IsInf(x, 1) {
		return "inf"
	}
	if math.IsInf(x, -1) {
		return "-inf"
	}
	if math.IsNaN(x) {
		if math.Signbit(x) {
			return "-nan"
		}
		return "nan"
	}
	return strconv.FormatFloat(x, 'g', 17, 64)
}
