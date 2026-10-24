// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/s_tan.c, utan.h, utan.tbl, dla.h.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func tanPolynomial(z float64) float64 {
	p := math.FMA(z, tangentTaylorUndecic, tangentTaylorNonic)
	p = math.FMA(z, p, tangentTaylorSeptic)
	p = math.FMA(z, p, tangentTaylorQuintic)
	return math.FMA(z, p, tangentTaylorCubic)
}

func tanTableResult(z float64, entry tanEntry, n int) float64 {
	z2 := float64(z * z)
	pz := math.FMA(float64(z*z2), math.FMA(z2, tangentTableQuintic, tangentTableCubic), z)
	fi, gi := entry.tangent, entry.cotangent
	if n != 0 {
		t2 := float64(pz*(fi+gi)) / (fi + pz)
		return -(gi - t2)
	}
	t2 := float64(pz*(gi+fi)) / (gi - pz)
	return fi + t2
}

func tanReduced(a, da float64, n int) float64 {
	ya, yya := a, da
	sy := float64(1)
	if a < 0 {
		ya, yya, sy = -a, -da, -1
	}
	if ya <= tangentTaylorLimit {
		a2 := float64(a * a)
		t2 := math.FMA(float64(a*a2), tanPolynomial(a2), da)
		if n != 0 {
			b, db := exactAdd(a, t2)
			c, dc := dividePair(1, 0, b, db)
			return -(c + dc)
		}
		return a + t2
	}
	i := int(math.FMA(lookupScale, ya, tangentTableIndexBias))
	entry := tanTable(i)
	z := (ya - entry.argument) + yya
	return float64(sy * tanTableResult(z, entry, n))
}

func Tan(x float64) float64 {
	if highWord(x)&exponentHighMask == exponentHighMask {
		return invalidTrig(x)
	}
	w := math.Abs(x)
	if w <= tangentIdentityLimit {
		return x
	}
	if w <= tangentTaylorLimit {
		x2 := float64(x * x)
		return math.FMA(float64(x*x2), tanPolynomial(x2), x)
	}
	if w <= tangentDirectTableLimit {
		i := int(math.FMA(lookupScale, w, tangentTableIndexBias))
		entry := tanTable(i)
		return math.Copysign(tanTableResult(w-entry.argument, entry, 0), x)
	}
	var n int
	var a, da float64
	if w <= tangentMediumReductionLimit {
		t := math.FMA(x, reciprocalPiHalf, integerRoundBias)
		xn := t - integerRoundBias
		t1 := math.FMA(-xn, piHalfTail, math.FMA(-xn, piHalfHead, x))
		n = lowWord(t) & 1
		if w <= tangentShortReductionLimit {
			a = math.FMA(-xn, piHalfCorrection, t1)
			da = math.FMA(-xn, piHalfCorrection, t1-a)
		} else {
			t = math.FMA(-xn, piHalfCorrectionHead, t1)
			da = math.FMA(-xn, piHalfCorrectionHead, t1-t)
			a = math.FMA(-xn, piHalfCorrectionTail, t)
			da = math.FMA(-xn, piHalfCorrectionTail, t-a) + da
			a, da = exactAdd(a, da)
		}
	} else {
		n, a, da = reduceLarge(x)
		n &= 1
		a, da = exactAdd(a, da)
	}
	return tanReduced(a, da, n)
}
