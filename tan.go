package libm

import "math"

const (
	tand3      float64 = 0x1.5555555555555p-2
	tand5      float64 = 0x1.11111111107c6p-3
	tand7      float64 = 0x1.ba1ba1cdb8745p-5
	tand9      float64 = 0x1.664ed49cfc666p-6
	tand11     float64 = 0x1.2385a3cf2e4eap-7
	tane0      float64 = 0x1.5555555554dbdp-2
	tane1      float64 = 0x1.11112e0a6b45fp-3
	tanmfftnhf float64 = -0x1.f000000000000p+3
	tang1      float64 = 0x1.b096c00000000p-27
	tang2      float64 = 0x1.f212d00000000p-5
	tang3      float64 = 0x1.92f1a00000000p-1
	tang4      float64 = 0x1.9000000000000p+4
	tang5      float64 = 0x1.7d78400000000p+26
	tangy2     float64 = 0x1.f212d00000000p-5
	tanmp1     float64 = 0x1.921fb58000000p+0
	tanmp2     float64 = -0x1.dde973c000000p-27
	tanmp3     float64 = -0x1.cb3b399d747f2p-55
	tanpp3     float64 = -0x1.cb3b398000000p-55
	tanpp4     float64 = -0x1.d747f23e32ed7p-83
	tanhpinv   float64 = 0x1.45f306dc9c883p-1
	tantoint   float64 = 0x1.8000000000000p+52
)

func tanPolynomial(z float64) float64 {
	p := math.FMA(z, tand11, tand9)
	p = math.FMA(z, p, tand7)
	p = math.FMA(z, p, tand5)
	return math.FMA(z, p, tand3)
}

func tanReduced(a, da float64, n int) float64 {
	ya, yya := a, da
	sy := float64(1)
	if a < 0 {
		ya, yya, sy = -a, -da, -1
	}
	if ya <= tangy2 {
		a2 := float64(a * a)
		t2 := math.FMA(float64(a*a2), tanPolynomial(a2), da)
		if n != 0 {
			b, db := exactAdd(a, t2)
			c, dc := dividePair(1, 0, b, db)
			return -float64(c + dc)
		}
		return float64(a + t2)
	}
	i := int(math.FMA(256, ya, tanmfftnhf))
	z := float64(float64(ya-tanTable(i*4)) + yya)
	z2 := float64(z * z)
	pz := math.FMA(float64(z*z2), math.FMA(z2, tane1, tane0), z)
	fi, gi := tanTable(i*4+1), tanTable(i*4+2)
	if n != 0 {
		t2 := float64(float64(pz*float64(fi+gi)) / float64(fi+pz))
		return float64(-sy * float64(gi-t2))
	}
	t2 := float64(float64(pz*float64(gi+fi)) / float64(gi-pz))
	return float64(sy * float64(fi+t2))
}

func Tan(x float64) float64 {
	if highWord(x)&0x7ff00000 == 0x7ff00000 {
		return invalidTrig(x)
	}
	w := math.Abs(x)
	if w <= tang1 {
		return x
	}
	if w <= tang2 {
		x2 := float64(x * x)
		return math.FMA(float64(x*x2), tanPolynomial(x2), x)
	}
	if w <= tang3 {
		i := int(math.FMA(256, w, tanmfftnhf))
		z := float64(w - tanTable(i*4))
		z2 := float64(z * z)
		pz := math.FMA(float64(z*z2), math.FMA(z2, tane1, tane0), z)
		fi, gi := tanTable(i*4+1), tanTable(i*4+2)
		t2 := float64(float64(pz*float64(gi+fi)) / float64(gi-pz))
		return math.Copysign(float64(fi+t2), x)
	}
	var n int
	var a, da float64
	if w <= tang5 {
		t := math.FMA(x, tanhpinv, tantoint)
		xn := float64(t - tantoint)
		t1 := math.FMA(-xn, tanmp2, math.FMA(-xn, tanmp1, x))
		n = lowWord(t) & 1
		if w <= tang4 {
			da = float64(xn * tanmp3)
			a = float64(t1 - da)
			da = float64(float64(t1-a) - da)
		} else {
			da = float64(xn * tanpp3)
			t = float64(t1 - da)
			da = float64(float64(t1-t) - da)
			t1 = float64(xn * tanpp4)
			a = float64(t - t1)
			da = float64(float64(float64(t-a)-t1) + da)
			a, da = exactAdd(a, da)
		}
	} else {
		n, a, da = reduceLarge(x)
		n &= 1
		a, da = exactAdd(a, da)
	}
	return tanReduced(a, da, n)
}

func invalidTrig(x float64) float64 {
	if math.IsNaN(x) {
		return math.Float64frombits(math.Float64bits(x) | 1<<51)
	}
	return math.Float64frombits(0xfff8000000000000)
}
