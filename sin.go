package libm

import "math"

const (
	sins1    float64 = -0x1.5555555555555p-3
	sins2    float64 = 0x1.1111111110ecep-7
	sins3    float64 = -0x1.a01a019db08b8p-13
	sins4    float64 = 0x1.71de27b9a7ed9p-19
	sins5    float64 = -0x1.addffc2fcdf59p-26
	sinbig   float64 = 0x1.8p45
	sinhp0   float64 = 0x1.921fb54442d18p0
	sinhp1   float64 = 0x1.1a62633145c07p-54
	sinmp1   float64 = 0x1.921fb58000000p0
	sinmp2   float64 = -0x1.dde973c000000p-27
	sinpp3   float64 = -0x1.cb3b398000000p-55
	sinpp4   float64 = -0x1.d747f23e32ed7p-83
	sinhpinv float64 = 0x1.45f306dc9c883p-1
	sintoint float64 = 0x1.8p52
	sinsn3   float64 = -1.66666666666664880952546298448555e-1
	sinsn5   float64 = 8.33333214285722277379541354343671e-3
	sincs2   float64 = 4.99999999999999999999950396842453e-1
	sincs4   float64 = -4.16666666666664434524222570944589e-2
	sincs6   float64 = 1.38888874007937613028114285595617e-3
)

func taylorSin(x, dx float64) float64 {
	xx := float64(x * x)
	p := math.FMA(sins5, xx, sins4)
	p = math.FMA(p, xx, sins3)
	p = math.FMA(p, xx, sins2)
	p = math.FMA(p, xx, sins1)
	t := math.FMA(p, x, -0.5*dx)
	t = math.FMA(t, xx, dx)
	return float64(x + t)
}

func sincosValues(u float64) (sn, ssn, cs, ccs float64) {
	k := lowWord(u) * 4
	return sincosTable(k), sincosTable(k + 1), sincosTable(k + 2), sincosTable(k + 3)
}

func doCos(x, dx float64) float64 {
	if x < 0 {
		dx = -dx
	}
	u := float64(sinbig + math.Abs(x))
	x = float64(float64(math.Abs(x)-float64(u-sinbig)) + dx)
	xx := float64(x * x)
	s := math.FMA(float64(x*xx), math.FMA(xx, sinsn5, sinsn3), x)
	c := float64(xx * math.FMA(xx, math.FMA(xx, sincs6, sincs4), sincs2))
	sn, ssn, cs, ccs := sincosValues(u)
	cor := math.FMA(-s, ssn, ccs)
	cor = math.FMA(-cs, c, cor)
	cor = math.FMA(-sn, s, cor)
	return float64(cs + cor)
}

func doSin(x, dx float64) float64 {
	old := x
	if math.Abs(x) < 0.126 {
		return taylorSin(x, dx)
	}
	if x <= 0 {
		dx = -dx
	}
	u := float64(sinbig + math.Abs(x))
	x = float64(math.Abs(x) - float64(u-sinbig))
	xx := float64(x * x)
	s := float64(x + math.FMA(float64(x*xx), math.FMA(xx, sinsn5, sinsn3), dx))
	c := math.FMA(x, dx, float64(xx*math.FMA(xx, math.FMA(xx, sincs6, sincs4), sincs2)))
	sn, ssn, cs, ccs := sincosValues(u)
	cor := math.FMA(s, ccs, ssn)
	cor = math.FMA(-sn, c, cor)
	cor = math.FMA(cs, s, cor)
	return math.Copysign(float64(sn+cor), old)
}

func reduceSinCos(x float64) (n int, a, da float64) {
	t := math.FMA(x, sinhpinv, sintoint)
	xn := float64(t - sintoint)
	y := math.FMA(-xn, sinmp2, math.FMA(-xn, sinmp1, x))
	n = lowWord(t) & 3
	t1 := float64(xn * sinpp3)
	t2 := float64(y - t1)
	db := float64(float64(y-t2) - t1)
	t1 = float64(xn * sinpp4)
	a = float64(t2 - t1)
	da = float64(db + float64(float64(t2-a)-t1))
	return
}

func doSinCos(a, da float64, n int) float64 {
	var result float64
	if n&1 != 0 {
		result = doCos(a, da)
	} else {
		result = doSin(a, da)
	}
	if n&2 != 0 {
		return -result
	}
	return result
}

func Sin(x float64) float64 {
	k := highWord(x) & 0x7fffffff
	if k < 0x3e500000 {
		return x
	}
	if k < 0x3feb6000 {
		return doSin(x, 0)
	}
	if k < 0x400368fd {
		return math.Copysign(doCos(float64(sinhp0-math.Abs(x)), sinhp1), x)
	}
	if k < 0x419921fb {
		n, a, da := reduceSinCos(x)
		return doSinCos(a, da, n)
	}
	if k < 0x7ff00000 {
		n, a, da := reduceLarge(x)
		return doSinCos(a, da, n)
	}
	return invalidTrig(x)
}

func Cos(x float64) float64 {
	k := highWord(x) & 0x7fffffff
	if k < 0x3e400000 {
		return 1
	}
	if k < 0x3feb6000 {
		return doCos(x, 0)
	}
	if k < 0x400368fd {
		y := float64(sinhp0 - math.Abs(x))
		a := float64(y + sinhp1)
		da := float64(float64(y-a) + sinhp1)
		return doSin(a, da)
	}
	if k < 0x419921fb {
		n, a, da := reduceSinCos(x)
		return doSinCos(a, da, n+1)
	}
	if k < 0x7ff00000 {
		n, a, da := reduceLarge(x)
		return doSinCos(a, da, n+1)
	}
	return invalidTrig(x)
}
