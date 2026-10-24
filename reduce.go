// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/branred.c, branred.h, dla.h, mydefs.h.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func reducePart(x float64) (sum, b, bb float64) {
	k := ((highWord(x) >> 20) & 2047) - largeReductionExponentBias
	k /= base24Digits
	if k < 0 {
		k = 0
	}
	gor := math.Float64frombits(uint64(largeReductionScaleHighWord-((k*base24Digits)<<20)) << 32)
	var r [6]float64
	for i := range r {
		r[i] = float64(float64(x*twoOverPi(k+i)) * gor)
		gor = float64(gor * base24FractionScale)
	}
	for i := 0; i < 3; i++ {
		s := (r[i] + integerRoundBias) - integerRoundBias
		sum += s
		r[i] -= s
	}
	t := float64(0)
	for i := 0; i < 6; i++ {
		t += r[5-i]
	}
	bb = (((((r[0] - t) + r[1]) + r[2]) + r[3]) + r[4]) + r[5]
	s := (t + integerRoundBias) - integerRoundBias
	sum += s
	t -= s
	b = t + bb
	bb = (t - b) + bb
	s = (sum + quadrantRoundBias) - quadrantRoundBias
	sum -= s
	return
}

func reduceLarge(x float64) (n int, a, aa float64) {
	x = float64(x * largeReductionInputScale)
	t := float64(x * mantissaSplitFactor)
	x1 := t - (t - x)
	x2 := x - x1
	sum1, b1, bb1 := reducePart(x1)
	sum2, b2, bb2 := reducePart(x2)
	sum := sum1 + sum2
	b, bb := exactAdd(b1, b2)
	if b > 0.5 {
		b -= 1
		sum += 1
	} else if b < -0.5 {
		b += 1
		sum -= 1
	}
	s := b + ((bb + bb1) + bb2)
	t = ((b - s) + bb) + (bb1 + bb2)
	b = float64(s * mantissaSplitFactor)
	t1 := b - (b - s)
	t2 := s - t1
	b = float64(s * piHalf)
	bb = (float64(t1*piHalfHead) - b) + float64(t1*branredPiHalfTail)
	bb += float64(t2 * piHalfHead)
	low := float64(t2*branredPiHalfTail) + float64(s*piHalfLow)
	low += float64(t * piHalf)
	bb += low
	a = b + bb
	aa = (b - a) + bb
	n = int(sum) & 3
	return
}
