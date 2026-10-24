// Copyright (C) 2001-2024 Free Software Foundation, Inc.
// IBM Accurate Mathematical Library, written by International Business Machines Corp.
// Translated to Go from glibc 2.39 sysdeps/ieee754/dbl-64/usncs.h, s_sin.c, utan.h, uasncs.h, atnat.h, atnat2.h, branred.h, branred.c, dla.h, e_asin.c, e_atan2.c.
// Modified by acneism on 2026-10-24: see NOTICE and docs/porting.md.
// SPDX-License-Identifier: LGPL-2.1-or-later

package libm

const (
	piHalf                      float64 = 0x1.921fb54442d18p0
	piHalfLow                   float64 = 0x1.1a62633145c07p-54
	piHalfHead                  float64 = 0x1.921fb58000000p0
	piHalfTail                  float64 = -0x1.dde973c000000p-27
	piHalfCorrection            float64 = -0x1.cb3b399d747f2p-55
	piHalfCorrectionHead        float64 = -0x1.cb3b398000000p-55
	piHalfCorrectionTail        float64 = -0x1.d747f23e32ed7p-83
	branredPiHalfTail           float64 = -0x1.dde9740000000p-27
	reciprocalPiHalf            float64 = 0x1.45f306dc9c883p-1
	pi                          float64 = 0x1.921fb54442d18p1
	piLow                       float64 = 0x1.1a62633145c07p-53
	piQuarter                   float64 = 0x1.921fb54442d18p-1
	piThreeQuarter              float64 = 0x1.2d97c7f3321d2p1
	integerRoundBias            float64 = 0x1.8p52
	quadrantRoundBias           float64 = 0x1.8p54
	nearestIntegerBias          float64 = 0x1p52
	mantissaSplitFactor         float64 = 134217729
	lookupScale                 float64 = 256
	sineTableRoundBias          float64 = 0x1.8p45
	sineTaylorCubic             float64 = -0x1.5555555555555p-3
	sineTaylorQuintic           float64 = 0x1.1111111110ecep-7
	sineTaylorSeptic            float64 = -0x1.a01a019db08b8p-13
	sineTaylorNonic             float64 = 0x1.71de27b9a7ed9p-19
	sineTaylorUndecic           float64 = -0x1.addffc2fcdf59p-26
	sineTableCubic              float64 = -1.66666666666664880952546298448555e-1
	sineTableQuintic            float64 = 8.33333214285722277379541354343671e-3
	cosineTableQuadratic        float64 = 4.99999999999999999999950396842453e-1
	cosineTableQuartic          float64 = -4.16666666666664434524222570944589e-2
	cosineTableSextic           float64 = 1.38888874007937613028114285595617e-3
	sineTaylorLimit             float64 = 0.126
	tangentTaylorCubic          float64 = 0x1.5555555555555p-2
	tangentTaylorQuintic        float64 = 0x1.11111111107c6p-3
	tangentTaylorSeptic         float64 = 0x1.ba1ba1cdb8745p-5
	tangentTaylorNonic          float64 = 0x1.664ed49cfc666p-6
	tangentTaylorUndecic        float64 = 0x1.2385a3cf2e4eap-7
	tangentTableCubic           float64 = 0x1.5555555554dbdp-2
	tangentTableQuintic         float64 = 0x1.11112e0a6b45fp-3
	tangentTableIndexBias       float64 = -15.5
	tangentIdentityLimit        float64 = 0x1.b096c00000000p-27
	tangentTaylorLimit          float64 = 0x1.f212d00000000p-5
	tangentDirectTableLimit     float64 = 0x1.92f1a00000000p-1
	tangentShortReductionLimit  float64 = 25
	tangentMediumReductionLimit float64 = 1e8
	asinTaylorCubic             float64 = 1.66666666666664110590506577996662e-1
	asinTaylorQuintic           float64 = 7.50000000026122686814431784722623e-2
	asinTaylorSeptic            float64 = 4.46428561421059750978517350006940e-2
	asinTaylorNonic             float64 = 3.03821268582119319911193410625235e-2
	asinTaylorUndecic           float64 = 2.23551211026525610742786300334557e-2
	asinTaylorTredecic          float64 = 1.81382903404565056280372531963613e-2
	asinRootRoundBias           float64 = 0x1p24
	acosRootSplitFactor         float64 = 0x1p27
	inverseRootConstant         float64 = 9.99999999859990725855365213134618e-1
	inverseRootLinear           float64 = 4.99999999495955425917856814202739e-1
	inverseRootQuadratic        float64 = 3.75017500867345182581453026130850e-1
	inverseRootCubic            float64 = 3.12523626554518656309172508769531e-1
	atanTaylorCubic             float64 = -0x1.5555555555555p-2
	atanTaylorQuintic           float64 = 0x1.99999999997fdp-3
	atanTaylorSeptic            float64 = -0x1.24924923f7603p-3
	atanTaylorNonic             float64 = 0x1.c71c6e5129a3bp-4
	atanTaylorUndecic           float64 = -0x1.7458022b13c25p-4
	atanTaylorTredecic          float64 = 0x1.375f08b31cbcep-4
	atanIdentityLimit           float64 = 0x1.bb67a00000000p-27
	atanTaylorLimit             float64 = 0.0625
	atanReciprocalLimit         float64 = 1
	atanReciprocalTableLimit    float64 = 16
	atanSaturationLimit         float64 = 0x1.49ff200000000p52
	atan2ScaleUp                float64 = 0x1p500
	atan2ScaleDown              float64 = 0x1p-500
	largeReductionInputScale    float64 = 0x1p-600
	base24FractionScale         float64 = 0x1p-24
)

const (
	absoluteHighMask                     = 0x7fffffff
	exponentHighMask                     = 0x7ff00000
	fractionHighMask                     = 0x000fffff
	sineIdentityHighWord                 = (1023 - 26) << 20
	cosineUnityHighWord                  = (1023 - 27) << 20
	sincosDirectTableHighWord            = 0x3feb6000
	sincosComplementHighWord             = 0x400368fd
	sincosMediumReductionHighWord        = 0x419921fb
	acosUnityHighWord                    = 0x3c880000
	asinTaylorHighWord                   = (1023 - 3) << 20
	asinQuarterHighWord                  = (1023 - 2) << 20
	asinHalfHighWord                     = (1023 - 1) << 20
	asinThreeQuarterHighWord             = 0x3fe80000
	asinFirstUpperHighWord               = 0x3fed8000
	asinSecondUpperHighWord              = 0x3fee8000
	asinTableUpperHighWord               = 0x3fef0000
	domainUnitHighWord                   = 1023 << 20
	atan2ExponentGapLimit                = 57 << 20
	asinQuarterSegmentOffset             = 352
	asinHalfSegmentOffset                = 1056
	asinThreeQuarterSegmentOffset        = 992
	asinFirstUpperSegmentOffset          = 884
	asinSecondUpperSegmentOffset         = 768
	atanTableFirstIndex                  = 16
	largeReductionScaleHighWord          = 0x63f00000
	largeReductionExponentBias           = 450
	base24Digits                         = 24
	quietNaNBit                   uint64 = 1 << 51
	positiveQuietNaNBits          uint64 = 0x7ff8000000000000
	negativeQuietNaNBits          uint64 = 0xfff8000000000000
)
