package grid

import (
	"math"
	"strconv"
)

const Count = 100000
const Functions = 7

func Name(function int) string {
	return [Functions]string{"sin", "cos", "tan", "asin", "acos", "atan", "atan2"}[function]
}

func mix(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

func unit(x uint64) float64 {
	return float64(x>>11) * 0x1p-53
}

func signed(x float64, bits uint64) float64 {
	if bits&1 != 0 {
		return -x
	}
	return x
}

func special(i uint64) float64 {
	return math.Float64frombits([12]uint64{0, 1 << 63, 1, 1<<63 | 1, 0x000fffffffffffff, 0x0010000000000000, 0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000000, 0xfff8000000000000, 0x3ff0000000000000, 0xbff0000000000000}[i%12])
}

func Input(function, index int) (x, y float64) {
	a := mix(uint64(index) + uint64(function)*0x100000001b3)
	b := mix(a)
	c := mix(b)
	mode := index % 8
	if index < 144 {
		return special(uint64(index)), special(uint64(index / 12))
	}
	y = math.Float64frombits(b)
	if function == 6 {
		switch mode {
		case 0:
			return float64(unit(a)*64 - 32), float64(unit(b)*64 - 32)
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
			return signed(float64(1+unit(a)*15), b), signed(float64(1+unit(b)*15), c)
		case 6:
			return math.Ldexp(signed(float64(1+unit(a)), b), int(c%120)-60), signed(float64(1+unit(b)), a)
		default:
			return math.Float64frombits(a & 0x800fffffffffffff), math.Float64frombits(b & 0x800fffffffffffff)
		}
	}
	if function == 3 || function == 4 {
		switch mode {
		case 0, 7:
			x = float64(unit(a)*2 - 1)
		case 1:
			x = signed(float64(1-math.Ldexp(unit(a), -int(b%52))), c)
		case 2:
			x = math.Float64frombits(a)
		case 3:
			x = signed(float64(float64(a%128)/128+math.Ldexp(unit(b), -12)), c)
		case 4:
			x = [8]float64{0.125, 0.25, 0.5, 0.75, 0.921875, 0.953125, 0.96875, 1}[a%8]
			x = math.Nextafter(x, signed(math.Inf(1), b))
			x = signed(x, c)
		case 5:
			x = math.Float64frombits(a & 0x800fffffffffffff)
		case 6:
			x = signed(math.Nextafter(1, 2), a)
		}
		return x, y
	}
	switch mode {
	case 0:
		x = float64(unit(a)*64 - 32)
	case 1:
		x = float64(unit(a)*2e8 - 1e8)
	case 2:
		x = math.Float64frombits(a)
	case 3:
		x = math.Ldexp(signed(float64(1+unit(a)), b), int(c%2098)-1074)
	case 4:
		x = float64(float64(int64(a%200000001)-100000000) * (math.Pi / 2))
		x = math.Nextafter(x, signed(math.Inf(1), b))
	case 5:
		x = signed(float64(float64(a%256)/256+math.Ldexp(unit(b), -16)), c)
	case 6:
		x = [12]float64{0x1p-27, 0x1p-26, 0.126, 0.85546875, 2.426265, 105414350, 0.0608, 0.787, 25, 1e8, 0.0625, 16}[a%12]
		x = signed(math.Nextafter(x, signed(math.Inf(1), b)), c)
	case 7:
		x = signed(math.Ldexp(unit(a), -int(b%1074)), c)
	}
	return x, y
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
