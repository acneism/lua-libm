// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/dla.h, mydefs.h, s_sin.c, s_tan.c, e_asin.c, s_atan.c, e_atan2.c.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

import "math"

type sincosEntry struct{ sine, sineTail, cosine, cosineTail float64 }
type tanEntry struct{ argument, tangent, cotangent, tangentTail float64 }
type atanEntry struct{ argument, value, linear, quadratic, cubic, quartic, quintic float64 }

func sincosTable(index int) sincosEntry {
	entry := sincosTableBits[index]
	return sincosEntry{math.Float64frombits(entry.sine), math.Float64frombits(entry.sineTail), math.Float64frombits(entry.cosine), math.Float64frombits(entry.cosineTail)}
}

func tanTable(index int) tanEntry {
	entry := tanTableBits[index]
	return tanEntry{math.Float64frombits(entry.argument), math.Float64frombits(entry.tangent), math.Float64frombits(entry.cotangent), math.Float64frombits(entry.tangentTail)}
}

func atanTable(index int) atanEntry {
	entry := atanTableBits[index]
	return atanEntry{math.Float64frombits(entry.argument), math.Float64frombits(entry.value), math.Float64frombits(entry.linear), math.Float64frombits(entry.quadratic), math.Float64frombits(entry.cubic), math.Float64frombits(entry.quartic), math.Float64frombits(entry.quintic)}
}

func asinTable(index int) float64  { return math.Float64frombits(asinTableBits[index]) }
func rootTable(index int) float64  { return math.Float64frombits(rootTableBits[index]) }
func powerTable(index int) float64 { return math.Float64frombits(powerTableBits[index]) }
func twoOverPi(index int) float64  { return math.Float64frombits(twoOverPiBits[index]) }

func highWord(x float64) int     { return int(math.Float64bits(x) >> 32) }
func lowWord(x float64) int      { return int(uint32(math.Float64bits(x))) }
func quietNaN(x float64) float64 { return math.Float64frombits(math.Float64bits(x) | quietNaNBit) }

func invalidTrig(x float64) float64 {
	if math.IsNaN(x) {
		return quietNaN(x)
	}
	return math.Float64frombits(negativeQuietNaNBits)
}

func exactAdd(x, y float64) (z, zz float64) {
	z = x + y
	if math.Abs(x) > math.Abs(y) {
		zz = (x - z) + y
	} else {
		zz = (y - z) + x
	}
	return
}

func exactSub(x, y float64) (z, zz float64) {
	z = x - y
	if math.Abs(x) > math.Abs(y) {
		zz = (x - z) - y
	} else {
		zz = x - (y + z)
	}
	return
}

func dividePair(x, xx, y, yy float64) (z, zz float64) {
	c := x / y
	u := float64(c * y)
	uu := math.FMA(c, y, -u)
	cc := math.FMA(-c, yy, ((x-u)-uu)+xx) / y
	z = c + cc
	zz = (c - z) + cc
	return
}
