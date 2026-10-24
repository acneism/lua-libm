package libm

import "math"

const (
	asinhp0 float64 = 0x1.921fb54442d18p+0
	asinhp1 float64 = 0x1.1a62633145c07p-54
	asinf1  float64 = 1.66666666666664110590506577996662e-01
	asinf2  float64 = 7.50000000026122686814431784722623e-02
	asinf3  float64 = 4.46428561421059750978517350006940e-02
	asinf4  float64 = 3.03821268582119319911193410625235e-02
	asinf5  float64 = 2.23551211026525610742786300334557e-02
	asinf6  float64 = 1.81382903404565056280372531963613e-02
	asint24 float64 = 16777216.0
	asint27 float64 = 134217728.0
	asinrt0 float64 = 9.99999999859990725855365213134618e-01
	asinrt1 float64 = 4.99999999495955425917856814202739e-01
	asinrt2 float64 = 3.75017500867345182581453026130850e-01
	asinrt3 float64 = 3.12523626554518656309172508769531e-01
)

func asinPolynomial(z float64) float64 {
	p := math.FMA(asinf6, z, asinf5)
	p = math.FMA(p, z, asinf4)
	p = math.FMA(p, z, asinf3)
	p = math.FMA(p, z, asinf2)
	return math.FMA(p, z, asinf1)
}

func asinIndex(k int) (n, degree int) {
	switch {
	case k < 0x3fe00000:
		if k < 0x3fd00000 {
			return 11 * ((k & 0xfffff) >> 15), 6
		}
		return 11*((k&0xfffff)>>14) + 352, 6
	case k < 0x3fe80000:
		return 1056 + ((k&0xfe000)>>11)*3, 7
	case k < 0x3fed8000:
		return 992 + ((k&0xfe000)>>13)*13, 8
	case k < 0x3fee8000:
		return 884 + ((k&0xfe000)>>13)*14, 9
	default:
		return 768 + ((k&0xfe000)>>13)*15, 10
	}
}

func asinApprox(x float64, k int) (base, correction float64) {
	n, degree := asinIndex(k)
	xx := float64(math.Abs(x) - asinTable(n))
	p := asinTable(n + degree)
	for j := degree - 1; j >= 2; j-- {
		p = math.FMA(xx, p, asinTable(n+j))
	}
	p = math.FMA(float64(xx*xx), p, asinTable(n+degree+1))
	correction = math.FMA(asinTable(n+1), xx, p)
	base = asinTable(n + degree + 2)
	return
}

func asinRoot(z float64, acos bool) (y, cc float64) {
	k := highWord(z)
	t := float64(rootTable((k&0x1fffff)>>14) * powerTable(511-(k>>21)))
	r := math.FMA(-float64(t*t), z, 1)
	t = float64(t * math.FMA(r, math.FMA(r, math.FMA(r, asinrt3, asinrt2), asinrt1), asinrt0))
	c := float64(t * z)
	factor := math.FMA(-float64(0.5*t), c, 1.5)
	if acos {
		y = math.FMA(-asint27, c, math.FMA(asint27, c, c))
	} else {
		y = float64(float64(c+asint24) - asint24)
	}
	cc = math.FMA(-y, y, z) / math.FMA(factor, c, y)
	return
}

func Asin(x float64) float64 {
	k := highWord(x) & 0x7fffffff
	switch {
	case k < 0x3e500000:
		return x
	case k < 0x3fc00000:
		xx := float64(x * x)
		return math.FMA(asinPolynomial(xx), float64(xx*x), x)
	case k < 0x3fef0000:
		base, correction := asinApprox(x, k)
		return math.Copysign(float64(base+correction), x)
	case k < 0x3ff00000:
		z := float64(0.5 * float64(1-math.Abs(x)))
		y, cc := asinRoot(z, false)
		p := float64(asinPolynomial(z) * z)
		cor := math.FMA(-float64(2*float64(y+cc)), p, math.FMA(-2, cc, asinhp1))
		res := float64(math.FMA(-2, y, asinhp0) + cor)
		return math.Copysign(res, x)
	case k == 0x3ff00000 && lowWord(x) == 0:
		return math.Copysign(asinhp0, x)
	default:
		return asinInvalid(x)
	}
}

func Acos(x float64) float64 {
	k := highWord(x) & 0x7fffffff
	switch {
	case k < 0x3c880000:
		return asinhp0
	case k < 0x3fc00000:
		xx := float64(x * x)
		r := float64(asinhp0 - x)
		cor := float64(float64(float64(asinhp0-r)-x) + asinhp1)
		return float64(r + math.FMA(-asinPolynomial(xx), float64(xx*x), cor))
	case k < 0x3fef0000:
		base, correction := asinApprox(x, k)
		if x > 0 {
			return float64(float64(asinhp0-base) + float64(asinhp1-correction))
		}
		return float64(float64(asinhp0+base) + float64(asinhp1+correction))
	case k < 0x3ff00000:
		z := float64(0.5 * float64(1-math.Abs(x)))
		y, cc := asinRoot(z, true)
		p := float64(asinPolynomial(z) * z)
		product := float64(p * float64(y+cc))
		if x < 0 {
			cor := float64(float64(asinhp1-cc) - product)
			res := float64(float64(asinhp0-y) + cor)
			return float64(res + res)
		}
		res := float64(y + float64(cc+product))
		return float64(res + res)
	case k == 0x3ff00000 && lowWord(x) == 0:
		if x > 0 {
			return 0
		}
		return 2 * asinhp0
	default:
		return asinInvalid(x)
	}
}

func asinInvalid(x float64) float64 {
	if math.IsNaN(x) {
		return math.Float64frombits(math.Float64bits(x) | 1<<51)
	}
	return math.Float64frombits(0x7ff8000000000000)
}
