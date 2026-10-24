// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/e_asin.c, uasncs.h, asincos.tbl, root.tbl, powtwo.tbl, mydefs.h.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

func asinPolynomial(z float64) float64 {
	p := math.FMA(asinTaylorTredecic, z, asinTaylorUndecic)
	p = math.FMA(p, z, asinTaylorNonic)
	p = math.FMA(p, z, asinTaylorSeptic)
	p = math.FMA(p, z, asinTaylorQuintic)
	return math.FMA(p, z, asinTaylorCubic)
}

func asinIndex(k int) (n, degree int) {
	switch {
	case k < asinHalfHighWord:
		if k < asinQuarterHighWord {
			return 11 * ((k & fractionHighMask) >> 15), 6
		}
		return 11*((k&fractionHighMask)>>14) + asinQuarterSegmentOffset, 6
	case k < asinThreeQuarterHighWord:
		return asinHalfSegmentOffset + ((k&0xfe000)>>11)*3, 7
	case k < asinFirstUpperHighWord:
		return asinThreeQuarterSegmentOffset + ((k&0xfe000)>>13)*13, 8
	case k < asinSecondUpperHighWord:
		return asinFirstUpperSegmentOffset + ((k&0xfe000)>>13)*14, 9
	default:
		return asinSecondUpperSegmentOffset + ((k&0xfe000)>>13)*15, 10
	}
}

func asinApprox(x float64, k int) (base, correction float64) {
	n, degree := asinIndex(k)
	xx := math.Abs(x) - asinTable(n)
	p := asinTable(n + degree)
	for j := degree - 1; j >= 2; j-- {
		p = math.FMA(xx, p, asinTable(n+j))
	}
	p = math.FMA(float64(xx*xx), p, asinTable(n+degree+1))
	correction = math.FMA(asinTable(n+1), xx, p)
	base = asinTable(n + degree + 2)
	return
}

func rootRefinement(z float64) (c, factor float64) {
	k := highWord(z)
	t := float64(rootTable((k&0x1fffff)>>14) * powerTable(511-(k>>21)))
	r := math.FMA(-float64(t*t), z, 1)
	t = float64(t * math.FMA(r, math.FMA(r, math.FMA(r, inverseRootCubic, inverseRootQuadratic), inverseRootLinear), inverseRootConstant))
	c = float64(t * z)
	factor = math.FMA(-float64(0.5*t), c, 1.5)
	return
}

func asinRoot(z float64) (y, cc float64) {
	c, factor := rootRefinement(z)
	y = (c + asinRootRoundBias) - asinRootRoundBias
	cc = math.FMA(-y, y, z) / math.FMA(factor, c, y)
	return
}

func acosRoot(z float64) (y, cc float64) {
	c, factor := rootRefinement(z)
	y = math.FMA(-acosRootSplitFactor, c, math.FMA(acosRootSplitFactor, c, c))
	cc = math.FMA(-y, y, z) / math.FMA(factor, c, y)
	return
}

func Asin(x float64) float64 {
	k := highWord(x) & absoluteHighMask
	switch {
	case k < sineIdentityHighWord:
		return x
	case k < asinTaylorHighWord:
		xx := float64(x * x)
		return math.FMA(asinPolynomial(xx), float64(xx*x), x)
	case k < asinTableUpperHighWord:
		base, correction := asinApprox(x, k)
		return math.Copysign(base+correction, x)
	case k < domainUnitHighWord:
		z := float64(0.5 * (1 - math.Abs(x)))
		y, cc := asinRoot(z)
		p := float64(asinPolynomial(z) * z)
		cor := math.FMA(-float64(2*(y+cc)), p, math.FMA(-2, cc, piHalfLow))
		res := math.FMA(-2, y, piHalf) + cor
		return math.Copysign(res, x)
	case k == domainUnitHighWord && lowWord(x) == 0:
		return math.Copysign(piHalf, x)
	default:
		return asinInvalid(x)
	}
}

func Acos(x float64) float64 {
	k := highWord(x) & absoluteHighMask
	switch {
	case k < acosUnityHighWord:
		return piHalf
	case k < asinTaylorHighWord:
		xx := float64(x * x)
		r := piHalf - x
		cor := ((piHalf - r) - x) + piHalfLow
		return r + math.FMA(-asinPolynomial(xx), float64(xx*x), cor)
	case k < asinTableUpperHighWord:
		base, correction := asinApprox(x, k)
		if x > 0 {
			return (piHalf - base) + (piHalfLow - correction)
		}
		return (piHalf + base) + (piHalfLow + correction)
	case k < domainUnitHighWord:
		z := float64(0.5 * (1 - math.Abs(x)))
		y, cc := acosRoot(z)
		p := float64(asinPolynomial(z) * z)
		product := float64(p * (y + cc))
		if x < 0 {
			cor := (piHalfLow - cc) - product
			res := (piHalf - y) + cor
			return res + res
		}
		res := y + (cc + product)
		return res + res
	case k == domainUnitHighWord && lowWord(x) == 0:
		if x > 0 {
			return 0
		}
		return float64(2 * piHalf)
	default:
		return asinInvalid(x)
	}
}

func asinInvalid(x float64) float64 {
	if math.IsNaN(x) {
		return quietNaN(x)
	}
	return math.Float64frombits(positiveQuietNaNBits)
}
