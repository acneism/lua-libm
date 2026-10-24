// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/s_atan.c, atnat.h, uatan.tbl, dla.h, mydefs.h.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func atanPolynomial(z float64) float64 {
	p := math.FMA(z, atanTaylorTredecic, atanTaylorUndecic)
	p = math.FMA(z, p, atanTaylorNonic)
	p = math.FMA(z, p, atanTaylorSeptic)
	p = math.FMA(z, p, atanTaylorQuintic)
	return math.FMA(z, p, atanTaylorCubic)
}

func atanTableHigherPolynomial(entry atanEntry, z float64) float64 {
	p := math.FMA(z, entry.quintic, entry.quartic)
	p = math.FMA(z, p, entry.cubic)
	return math.FMA(z, p, entry.quadratic)
}

func atanTablePolynomial(entry atanEntry, z float64) float64 {
	return math.FMA(z, atanTableHigherPolynomial(entry, z), entry.linear)
}

func atanIndex(u float64) int {
	return int(math.FMA(lookupScale, u, nearestIntegerBias)-nearestIntegerBias) - atanTableFirstIndex
}

func atanReciprocal(u float64) (w, correction float64) {
	w = 1 / u
	t1 := float64(w * u)
	t2 := math.FMA(w, u, -t1)
	correction = (1 - t1) - t2
	return
}

func Atan(x float64) float64 {
	if math.IsNaN(x) {
		return quietNaN(x)
	}
	u := math.Abs(x)
	switch {
	case u < atanIdentityLimit:
		return x
	case u < atanTaylorLimit:
		v := float64(x * x)
		return math.FMA(float64(x*v), atanPolynomial(v), x)
	case u < atanReciprocalLimit:
		entry := atanTable(atanIndex(u))
		z := u - entry.argument
		return math.Copysign(math.FMA(z, atanTablePolynomial(entry, z), entry.value), x)
	case u < atanReciprocalTableLimit:
		w, correction := atanReciprocal(u)
		entry := atanTable(atanIndex(w))
		z := math.FMA(correction, w, w-entry.argument)
		yy := math.FMA(-z, atanTablePolynomial(entry, z), piHalfLow)
		t1 := piHalf - entry.value
		return math.Copysign(t1+yy, x)
	case u < atanSaturationLimit:
		w, correction := atanReciprocal(u)
		v := float64(w * w)
		t3, cor := exactSub(piHalf, w)
		yy := math.FMA(-correction, w, piHalfLow+cor)
		yy = math.FMA(-float64(w*v), atanPolynomial(v), yy)
		return math.Copysign(t3+yy, x)
	default:
		return math.Copysign(piHalf, x)
	}
}
