package libm

import "math"

const (
	reduceSplit float64 = 134217729
	reduceBig   float64 = 0x1.8p52
	reduceBig1  float64 = 0x1.8p54
	reduceHp0   float64 = 0x1.921fb54442d18p0
	reduceHp1   float64 = 0x1.1a62633145c07p-54
	reduceMp1   float64 = 0x1.921fb58000000p0
	reduceMp2   float64 = -0x1.dde9740000000p-27
)

func reducePart(x float64) (sum, b, bb float64) {
	k := ((highWord(x) >> 20) & 2047) - 450
	k /= 24
	if k < 0 {
		k = 0
	}
	gor := math.Float64frombits(uint64(0x63f00000-((k*24)<<20)) << 32)
	var r [6]float64
	for i := range r {
		r[i] = float64(float64(x*twoOverPi(k+i)) * gor)
		gor = float64(gor * 0x1p-24)
	}
	for i := 0; i < 3; i++ {
		s := float64(float64(r[i]+reduceBig) - reduceBig)
		sum = float64(sum + s)
		r[i] = float64(r[i] - s)
	}
	t := float64(0)
	for i := 0; i < 6; i++ {
		t = float64(t + r[5-i])
	}
	bb = float64(float64(float64(float64(float64(float64(r[0]-t)+r[1])+r[2])+r[3])+r[4]) + r[5])
	s := float64(float64(t+reduceBig) - reduceBig)
	sum = float64(sum + s)
	t = float64(t - s)
	b = float64(t + bb)
	bb = float64(float64(t-b) + bb)
	s = float64(float64(sum+reduceBig1) - reduceBig1)
	sum = float64(sum - s)
	return
}

func reduceLarge(x float64) (n int, a, aa float64) {
	x = float64(x * 0x1p-600)
	t := float64(x * reduceSplit)
	x1 := float64(t - float64(t-x))
	x2 := float64(x - x1)
	sum1, b1, bb1 := reducePart(x1)
	sum2, b2, bb2 := reducePart(x2)
	sum := float64(sum1 + sum2)
	b, bb := exactAdd(b1, b2)
	if b > 0.5 {
		b = float64(b - 1)
		sum = float64(sum + 1)
	} else if b < -0.5 {
		b = float64(b + 1)
		sum = float64(sum - 1)
	}
	s := float64(b + float64(float64(bb+bb1)+bb2))
	t = float64(float64(float64(b-s)+bb) + float64(bb1+bb2))
	b = float64(s * reduceSplit)
	t1 := float64(b - float64(b-s))
	t2 := float64(s - t1)
	b = float64(s * reduceHp0)
	bb = math.FMA(t1, reduceMp1, -b)
	bb = math.FMA(t1, reduceMp2, bb)
	bb = math.FMA(t2, reduceMp1, bb)
	low := math.FMA(t2, reduceMp2, float64(s*reduceHp1))
	low = math.FMA(t, reduceHp0, low)
	bb = float64(bb + low)
	a = float64(b + bb)
	aa = float64(float64(b-a) + bb)
	n = int(sum) & 3
	return
}
