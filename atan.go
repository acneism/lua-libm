package libm

import "math"

const (
	atand3  float64 = -0x1.5555555555555p-2
	atand5  float64 = 0x1.99999999997fdp-3
	atand7  float64 = -0x1.24924923f7603p-3
	atand9  float64 = 0x1.c71c6e5129a3bp-4
	atand11 float64 = -0x1.7458022b13c25p-4
	atand13 float64 = 0x1.375f08b31cbcep-4
)
const atanA float64 = 0x1.bb67a00000000p-27
const atanB float64 = 0x1.0000000000000p-4
const atanC float64 = 0x1.0000000000000p+0
const atanD float64 = 0x1.0000000000000p+4
const atanE float64 = 0x1.49ff200000000p+52
const atanHPI float64 = 0x1.921fb54442d18p+0
const atanHPI1 float64 = 0x1.1a62633145c07p-54

func atanPolynomial(z float64) float64 {
	p := math.FMA(z, atand13, atand11)
	p = math.FMA(z, p, atand9)
	p = math.FMA(z, p, atand7)
	p = math.FMA(z, p, atand5)
	return math.FMA(z, p, atand3)
}

func atanTablePolynomial(i int, z float64) float64 {
	p := math.FMA(z, atanTable(i*7+6), atanTable(i*7+5))
	p = math.FMA(z, p, atanTable(i*7+4))
	p = math.FMA(z, p, atanTable(i*7+3))
	return math.FMA(z, p, atanTable(i*7+2))
}

func atanIndex(u float64) int {
	return int(float64(math.FMA(256, u, 0x1p52)-0x1p52)) - 16
}

func Atan(x float64) float64 {
	if math.IsNaN(x) {
		return math.Float64frombits(math.Float64bits(x) | 1<<51)
	}
	u := math.Abs(x)
	switch {
	case u < atanA:
		return x
	case u < atanB:
		v := float64(x * x)
		return math.FMA(float64(x*v), atanPolynomial(v), x)
	case u < atanC:
		i := atanIndex(u)
		z := float64(u - atanTable(i*7))
		return math.Copysign(math.FMA(z, atanTablePolynomial(i, z), atanTable(i*7+1)), x)
	case u < atanD:
		w := float64(1 / u)
		t1 := float64(w * u)
		t2 := math.FMA(w, u, -t1)
		cor := float64(float64(1-t1) - t2)
		i := atanIndex(w)
		z := math.FMA(cor, w, float64(w-atanTable(i*7)))
		yy := math.FMA(-z, atanTablePolynomial(i, z), atanHPI1)
		t1 = float64(atanHPI - atanTable(i*7+1))
		return math.Copysign(float64(t1+yy), x)
	case u < atanE:
		w := float64(1 / u)
		v := float64(w * w)
		t1 := float64(w * u)
		t2 := math.FMA(w, u, -t1)
		correction := float64(float64(1-t1) - t2)
		t3, cor := exactSub(atanHPI, w)
		yy := math.FMA(-correction, w, float64(atanHPI1+cor))
		yy = math.FMA(-float64(w*v), atanPolynomial(v), yy)
		return math.Copysign(float64(t3+yy), x)
	default:
		return math.Copysign(atanHPI, x)
	}
}
