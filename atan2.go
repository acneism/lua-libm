// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/e_atan2.c, atnat2.h, uatan.tbl, dla.h, mydefs.h.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func atan2Special(y, x float64) (float64, bool) {
	if math.IsNaN(x) {
		return quietNaN(x), true
	}
	if math.IsNaN(y) {
		return quietNaN(y), true
	}
	if y == 0 {
		if math.Signbit(x) {
			return math.Copysign(pi, y), true
		}
		return math.Copysign(0, y), true
	}
	if x == 0 {
		return math.Copysign(piHalf, y), true
	}
	if math.IsInf(x, 0) {
		if math.IsInf(y, 0) {
			if x > 0 {
				return math.Copysign(piQuarter, y), true
			}
			return math.Copysign(piThreeQuarter, y), true
		}
		if x > 0 {
			return math.Copysign(0, y), true
		}
		return math.Copysign(pi, y), true
	}
	if math.IsInf(y, 0) {
		return math.Copysign(piHalf, y), true
	}
	return 0, false
}

func atan2Ratio(ax, ay float64) (u, du float64) {
	if ax < atan2ScaleDown || ay < atan2ScaleDown {
		ax = float64(ax * atan2ScaleUp)
		ay = float64(ay * atan2ScaleUp)
	}
	if ax > atan2ScaleUp || ay > atan2ScaleUp {
		ax = float64(ax * atan2ScaleDown)
		ay = float64(ay * atan2ScaleDown)
	}
	if ay < ax {
		u = ay / ax
		v := float64(ax * u)
		vv := math.FMA(ax, u, -v)
		du = ((ay - v) - vv) / ax
	} else {
		u = ax / ay
		v := float64(ay * u)
		vv := math.FMA(ay, u, -v)
		du = ((ax - v) - vv) / ay
	}
	return
}

func atan2Direct(u, du float64) float64 {
	if u < atanTaylorLimit {
		v := float64(u * u)
		zz := math.FMA(float64(u*v), atanPolynomial(v), du)
		return u + zz
	}
	entry := atanTable(atanIndex(u))
	t3 := u - entry.argument
	v, dv := exactAdd(t3, du)
	p := float64(float64(v*v) * atanTableHigherPolynomial(entry, v))
	zz := math.FMA(v, entry.linear, math.FMA(dv, entry.linear, p))
	return entry.value + zz
}

func atan2AddOffset(u, du, base, low float64) float64 {
	if u < atanTaylorLimit {
		v := float64(u * u)
		zz := float64(float64(u*v) * atanPolynomial(v))
		t2, cor := exactAdd(base, u)
		t3 := ((low + cor) + du) + zz
		return t2 + t3
	}
	entry := atanTable(atanIndex(u))
	v := (u - entry.argument) + du
	zz := math.FMA(v, atanTablePolynomial(entry, v), low)
	t1 := base + entry.value
	return t1 + zz
}

func atan2SubtractOffset(u, du, base, low float64) float64 {
	if u < atanTaylorLimit {
		v := float64(u * u)
		zz := float64(float64(u*v) * atanPolynomial(v))
		t2, cor := exactSub(base, u)
		t3 := ((low + cor) - du) - zz
		return t2 + t3
	}
	entry := atanTable(atanIndex(u))
	v := (u - entry.argument) + du
	zz := math.FMA(-v, atanTablePolynomial(entry, v), low)
	t1 := base - entry.value
	return t1 + zz
}

func Atan2(y, x float64) float64 {
	if result, ok := atan2Special(y, x); ok {
		return result
	}
	ax, ay := math.Abs(x), math.Abs(y)
	de := (highWord(y) & exponentHighMask) - (highWord(x) & exponentHighMask)
	if de >= atan2ExponentGapLimit {
		return math.Copysign(piHalf, y)
	}
	if de <= -atan2ExponentGapLimit {
		if x > 0 {
			return math.Copysign(ay/ax, y)
		}
		return math.Copysign(pi, y)
	}
	u, du := atan2Ratio(ax, ay)
	var result float64
	switch {
	case x > 0 && ay < ax:
		result = atan2Direct(u, du)
	case x > 0:
		result = atan2SubtractOffset(u, du, piHalf, piHalfLow)
	case ax < ay:
		result = atan2AddOffset(u, du, piHalf, piHalfLow)
	default:
		result = atan2SubtractOffset(u, du, pi, piLow)
	}
	return math.Copysign(result, y)
}
