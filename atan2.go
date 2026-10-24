package libm

import "math"

const (
	atan2Pi             float64 = 0x1.921fb54442d18p1
	atan2PiLow          float64 = 0x1.1a62633145c07p-53
	atan2HalfPi         float64 = 0x1.921fb54442d18p0
	atan2HalfPiLow      float64 = 0x1.1a62633145c07p-54
	atan2QuarterPi      float64 = 0x1.921fb54442d18p-1
	atan2ThreeQuarterPi float64 = 0x1.2d97c7f3321d2p1
)

func atan2Special(y, x float64) (float64, bool) {
	if math.IsNaN(x) {
		return math.Float64frombits(math.Float64bits(x) | 1<<51), true
	}
	if math.IsNaN(y) {
		return math.Float64frombits(math.Float64bits(y) | 1<<51), true
	}
	if y == 0 {
		if math.Signbit(x) {
			return math.Copysign(atan2Pi, y), true
		}
		return math.Copysign(0, y), true
	}
	if x == 0 {
		return math.Copysign(atan2HalfPi, y), true
	}
	if math.IsInf(x, 0) {
		if math.IsInf(y, 0) {
			if x > 0 {
				return math.Copysign(atan2QuarterPi, y), true
			}
			return math.Copysign(atan2ThreeQuarterPi, y), true
		}
		if x > 0 {
			return math.Copysign(0, y), true
		}
		return math.Copysign(atan2Pi, y), true
	}
	if math.IsInf(y, 0) {
		return math.Copysign(atan2HalfPi, y), true
	}
	return 0, false
}

func atan2Ratio(ax, ay float64) (u, du float64) {
	if ax < 0x1p-500 || ay < 0x1p-500 {
		ax = float64(ax * 0x1p500)
		ay = float64(ay * 0x1p500)
	}
	if ax > 0x1p500 || ay > 0x1p500 {
		ax = float64(ax * 0x1p-500)
		ay = float64(ay * 0x1p-500)
	}
	if ay < ax {
		u = float64(ay / ax)
		v := float64(ax * u)
		vv := math.FMA(ax, u, -v)
		du = float64(float64(float64(ay-v)-vv) / ax)
	} else {
		u = float64(ax / ay)
		v := float64(ay * u)
		vv := math.FMA(ay, u, -v)
		du = float64(float64(float64(ax-v)-vv) / ay)
	}
	return
}

func atan2Direct(u, du float64) float64 {
	if u < 0.0625 {
		v := float64(u * u)
		zz := math.FMA(float64(u*v), atanPolynomial(v), du)
		return float64(u + zz)
	}
	i := atanIndex(u)
	t3 := float64(u - atanTable(i*7))
	v, dv := exactAdd(t3, du)
	t1, t2 := atanTable(i*7+1), atanTable(i*7+2)
	p := math.FMA(v, atanTable(i*7+6), atanTable(i*7+5))
	p = math.FMA(v, p, atanTable(i*7+4))
	p = math.FMA(v, p, atanTable(i*7+3))
	p = float64(float64(v*v) * p)
	zz := math.FMA(v, t2, math.FMA(dv, t2, p))
	return float64(t1 + zz)
}

func atan2Offset(u, du, base, low float64, add bool) float64 {
	if u < 0.0625 {
		v := float64(u * u)
		zz := float64(float64(u*v) * atanPolynomial(v))
		if add {
			t2, cor := exactAdd(base, u)
			t3 := float64(float64(float64(low+cor)+du) + zz)
			return float64(t2 + t3)
		}
		t2, cor := exactSub(base, u)
		t3 := float64(float64(float64(low+cor)-du) - zz)
		return float64(t2 + t3)
	}
	i := atanIndex(u)
	v := float64(float64(u-atanTable(i*7)) + du)
	p := atanTablePolynomial(i, v)
	if add {
		zz := math.FMA(v, p, low)
		t1 := float64(base + atanTable(i*7+1))
		return float64(t1 + zz)
	}
	zz := math.FMA(-v, p, low)
	t1 := float64(base - atanTable(i*7+1))
	return float64(t1 + zz)
}

func Atan2(y, x float64) float64 {
	if result, ok := atan2Special(y, x); ok {
		return result
	}
	ax, ay := math.Abs(x), math.Abs(y)
	de := (highWord(y) & 0x7ff00000) - (highWord(x) & 0x7ff00000)
	if de >= 59768832 {
		return math.Copysign(atan2HalfPi, y)
	}
	if de <= -59768832 {
		if x > 0 {
			return math.Copysign(float64(ay/ax), y)
		}
		return math.Copysign(atan2Pi, y)
	}
	u, du := atan2Ratio(ax, ay)
	var result float64
	switch {
	case x > 0 && ay < ax:
		result = atan2Direct(u, du)
	case x > 0:
		result = atan2Offset(u, du, atan2HalfPi, atan2HalfPiLow, false)
	case ax < ay:
		result = atan2Offset(u, du, atan2HalfPi, atan2HalfPiLow, true)
	default:
		result = atan2Offset(u, du, atan2Pi, atan2PiLow, false)
	}
	return math.Copysign(result, y)
}
