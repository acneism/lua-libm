// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/s_sin.c, usncs.h, sincostab.c.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func taylorSin(x, dx float64) float64 {
	xx := float64(x * x)
	p := math.FMA(sineTaylorUndecic, xx, sineTaylorNonic)
	p = math.FMA(p, xx, sineTaylorSeptic)
	p = math.FMA(p, xx, sineTaylorQuintic)
	p = math.FMA(p, xx, sineTaylorCubic)
	t := math.FMA(p, x, float64(-0.5*dx))
	t = math.FMA(t, xx, dx)
	return x + t
}

func doCos(x, dx float64) float64 {
	if x < 0 {
		dx = -dx
	}
	u := sineTableRoundBias + math.Abs(x)
	x = (math.Abs(x) - (u - sineTableRoundBias)) + dx
	xx := float64(x * x)
	s := math.FMA(float64(x*xx), math.FMA(xx, sineTableQuintic, sineTableCubic), x)
	c := float64(xx * math.FMA(xx, math.FMA(xx, cosineTableSextic, cosineTableQuartic), cosineTableQuadratic))
	entry := sincosTable(lowWord(u))
	cor := math.FMA(-s, entry.sineTail, entry.cosineTail)
	cor = math.FMA(-entry.cosine, c, cor)
	cor = math.FMA(-entry.sine, s, cor)
	return entry.cosine + cor
}

func doSin(x, dx float64) float64 {
	old := x
	if math.Abs(x) < sineTaylorLimit {
		return taylorSin(x, dx)
	}
	if x <= 0 {
		dx = -dx
	}
	u := sineTableRoundBias + math.Abs(x)
	x = math.Abs(x) - (u - sineTableRoundBias)
	xx := float64(x * x)
	s := x + math.FMA(float64(x*xx), math.FMA(xx, sineTableQuintic, sineTableCubic), dx)
	c := math.FMA(x, dx, float64(xx*math.FMA(xx, math.FMA(xx, cosineTableSextic, cosineTableQuartic), cosineTableQuadratic)))
	entry := sincosTable(lowWord(u))
	cor := math.FMA(s, entry.cosineTail, entry.sineTail)
	cor = math.FMA(-entry.sine, c, cor)
	cor = math.FMA(entry.cosine, s, cor)
	return math.Copysign(entry.sine+cor, old)
}

func reduceSinCos(x float64) (n int, a, da float64) {
	t := math.FMA(x, reciprocalPiHalf, integerRoundBias)
	xn := t - integerRoundBias
	y := math.FMA(-xn, piHalfTail, math.FMA(-xn, piHalfHead, x))
	n = lowWord(t) & 3
	t2 := math.FMA(-xn, piHalfCorrectionHead, y)
	db := math.FMA(-xn, piHalfCorrectionHead, y-t2)
	a = math.FMA(-xn, piHalfCorrectionTail, t2)
	da = db + math.FMA(-xn, piHalfCorrectionTail, t2-a)
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
	k := highWord(x) & absoluteHighMask
	if k < sineIdentityHighWord {
		return x
	}
	if k < sincosDirectTableHighWord {
		return doSin(x, 0)
	}
	if k < sincosComplementHighWord {
		return math.Copysign(doCos(piHalf-math.Abs(x), piHalfLow), x)
	}
	if k < sincosMediumReductionHighWord {
		n, a, da := reduceSinCos(x)
		return doSinCos(a, da, n)
	}
	if k < exponentHighMask {
		n, a, da := reduceLarge(x)
		return doSinCos(a, da, n)
	}
	return invalidTrig(x)
}

func Cos(x float64) float64 {
	k := highWord(x) & absoluteHighMask
	if k < cosineUnityHighWord {
		return 1
	}
	if k < sincosDirectTableHighWord {
		return doCos(x, 0)
	}
	if k < sincosComplementHighWord {
		y := piHalf - math.Abs(x)
		a := y + piHalfLow
		da := (y - a) + piHalfLow
		return doSin(a, da)
	}
	if k < sincosMediumReductionHighWord {
		n, a, da := reduceSinCos(x)
		return doSinCos(a, da, n+1)
	}
	if k < exponentHighMask {
		n, a, da := reduceLarge(x)
		return doSinCos(a, da, n+1)
	}
	return invalidTrig(x)
}
