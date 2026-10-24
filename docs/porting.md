# Porting

The source version is [glibc 2.39](https://github.com/bminor/glibc/tree/glibc-2.39/sysdeps/ieee754/dbl-64). [sources.sha256](sources.sha256) records the SHA-256 of every original numerical source and header used by this port. The source bytes remain unchanged. Each translated Go file carries the upstream copyright, its own source list, the modification date, and the LGPL-2.1-or-later identifier. The headers and the generated-file marker are the only Go code comments.

| Go file | glibc sources | Translation changes |
| --- | --- | --- |
| arithmetic.go | dla.h, mydefs.h, s_sin.c, s_tan.c, e_asin.c, s_atan.c, e_atan2.c | Replace unions with Float64bits; share exact-add/subtract/divide operations and explicit quiet-NaN selection; decode typed table entries. |
| sin.go | s_sin.c, usncs.h, sincostab.c | Share the sine/cosine polynomial helpers; preserve table fields, reduction order, quadrant selection, and each compiler-selected FMA. |
| tan.go | s_tan.c, utan.h, utan.tbl, dla.h | Share the direct and reduced table formula; fuse both uses of the reduction products, including xn*mp3 near poles. |
| asin.go | e_asin.c, uasncs.h, asincos.tbl, root.tbl, powtwo.tbl, mydefs.h | Name the five segment offsets; share root refinement and expose separate asin/acos root functions; preserve the distinct final corrections. |
| atan.go | s_atan.c, atnat.h, uatan.tbl, dla.h, mydefs.h | Share reciprocal correction and table polynomials; preserve FMA boundaries across the original assignments. |
| atan2.go | e_atan2.c, atnat2.h, uatan.tbl, dla.h, mydefs.h | Share the atan coefficients; use separate additive and subtractive quadrant corrections; preserve special-value argument priority. |
| reduce.go | branred.c, branred.h, dla.h, mydefs.h | Share the two base-24 reduction blocks. Keep all products separately rounded because glibc explicitly compiles branred without FMA. |
| constants.go | usncs.h, s_sin.c, utan.h, uasncs.h, atnat.h, atnat2.h, branred.h, branred.c, dla.h, e_asin.c, e_atan2.c | Centralize shared coefficients, range thresholds, masks, offsets, and scaling values under descriptive names. |
| tables.go | sincostab.c, utan.tbl, uatan.tbl, asincos.tbl, root.tbl, powtwo.tbl, branred.h | Generate ordered uint64 literals with typed row layouts. Preserve every original binary64 bit. |

The public Go functions do not expose C errno, floating-point exception flags, or directed rounding modes. Explicit NaN selection follows the public libm entry points used by Redis. No numerical expression is algebraically reassociated during translation.

## Lookup tables

Tables are package-level arrays that the numerical code only reads. Each hexadecimal uint64 literal joins glibc's high and low 32-bit words in that order. The sine/cosine rows have named sine, sineTail, cosine, and cosineTail fields. The tangent rows have argument, tangent, cotangent, and tangentTail fields. The atan rows have argument, value, linear, quadratic, cubic, quartic, and quintic fields. Variable-length asin segments retain their upstream flat order and use named base offsets.

`internal/cmd/gentables` reads the unmodified upstream files, verifies every file against the manifest before parsing, validates table sizes, and emits gofmt-formatted tables. A missing file, changed source hash, malformed word pair, or wrong element count stops generation. Decimal root/power/base-24 values are converted with strconv.ParseFloat and Float64bits.

```text
go run ./internal/cmd/gentables -source work/source -manifest docs/sources.sha256 -output tables.go
```

Regenerating from the verified source bytes reproduces tables.go byte for byte. The generator and its tests are original Go tooling and do not need translated-source headers.

## Constants

Every numerical constant in constants.go is listed below. Slash-separated upstream names denote equal constants that now share one definition. Threshold descriptions denote literals in the corresponding C branch. IEEE 754 masks name the binary64 encodings used by those branches. Float constants are typed float64 so decimal values round to binary64 before evaluation. The branred mp2 remains separate because its bits differ from the sin/tan mp2.

| Go name | glibc name or expression | Go value |
| --- | --- | --- |
| `piHalf` | `usncs.h:hp0/hpi/HPI` | `0x1.921fb54442d18p0` |
| `piHalfLow` | `usncs.h:hp1/hpi1/HPI1` | `0x1.1a62633145c07p-54` |
| `piHalfHead` | `usncs.h:mp1` | `0x1.921fb58000000p0` |
| `piHalfTail` | `usncs.h:mp2` | `-0x1.dde973c000000p-27` |
| `piHalfCorrection` | `usncs.h:mp3` | `-0x1.cb3b399d747f2p-55` |
| `piHalfCorrectionHead` | `usncs.h:pp3` | `-0x1.cb3b398000000p-55` |
| `piHalfCorrectionTail` | `usncs.h:pp4` | `-0x1.d747f23e32ed7p-83` |
| `branredPiHalfTail` | `branred.h:mp2` | `-0x1.dde9740000000p-27` |
| `reciprocalPiHalf` | `usncs.h:hpinv` | `0x1.45f306dc9c883p-1` |
| `pi` | `atnat2.h:opi` | `0x1.921fb54442d18p1` |
| `piLow` | `atnat2.h:opi1` | `0x1.1a62633145c07p-53` |
| `piQuarter` | `atnat2.h:qpi` | `0x1.921fb54442d18p-1` |
| `piThreeQuarter` | `atnat2.h:tqpi` | `0x1.2d97c7f3321d2p1` |
| `integerRoundBias` | `usncs.h:toint/branred.h:big` | `0x1.8p52` |
| `quadrantRoundBias` | `branred.h:big1` | `0x1.8p54` |
| `nearestIntegerBias` | `s_atan.c:TWO52` | `0x1p52` |
| `mantissaSplitFactor` | `dla.h:CN` | `134217729` |
| `lookupScale` | `s_atan.c:256` | `256` |
| `sineTableRoundBias` | `usncs.h:big` | `0x1.8p45` |
| `sineTaylorCubic` | `usncs.h:s1` | `-0x1.5555555555555p-3` |
| `sineTaylorQuintic` | `usncs.h:s2` | `0x1.1111111110ecep-7` |
| `sineTaylorSeptic` | `usncs.h:s3` | `-0x1.a01a019db08b8p-13` |
| `sineTaylorNonic` | `usncs.h:s4` | `0x1.71de27b9a7ed9p-19` |
| `sineTaylorUndecic` | `usncs.h:s5` | `-0x1.addffc2fcdf59p-26` |
| `sineTableCubic` | `s_sin.c:sn3` | `-1.66666666666664880952546298448555e-1` |
| `sineTableQuintic` | `s_sin.c:sn5` | `8.33333214285722277379541354343671e-3` |
| `cosineTableQuadratic` | `s_sin.c:cs2` | `4.99999999999999999999950396842453e-1` |
| `cosineTableQuartic` | `s_sin.c:cs4` | `-4.16666666666664434524222570944589e-2` |
| `cosineTableSextic` | `s_sin.c:cs6` | `1.38888874007937613028114285595617e-3` |
| `sineTaylorLimit` | `s_sin.c:do_sin_threshold` | `0.126` |
| `tangentTaylorCubic` | `utan.h:d3` | `0x1.5555555555555p-2` |
| `tangentTaylorQuintic` | `utan.h:d5` | `0x1.11111111107c6p-3` |
| `tangentTaylorSeptic` | `utan.h:d7` | `0x1.ba1ba1cdb8745p-5` |
| `tangentTaylorNonic` | `utan.h:d9` | `0x1.664ed49cfc666p-6` |
| `tangentTaylorUndecic` | `utan.h:d11` | `0x1.2385a3cf2e4eap-7` |
| `tangentTableCubic` | `utan.h:e0` | `0x1.5555555554dbdp-2` |
| `tangentTableQuintic` | `utan.h:e1` | `0x1.11112e0a6b45fp-3` |
| `tangentTableIndexBias` | `utan.h:mfftnhf` | `-15.5` |
| `tangentIdentityLimit` | `utan.h:g1` | `0x1.b096c00000000p-27` |
| `tangentTaylorLimit` | `utan.h:g2/gy2` | `0x1.f212d00000000p-5` |
| `tangentDirectTableLimit` | `utan.h:g3` | `0x1.92f1a00000000p-1` |
| `tangentShortReductionLimit` | `utan.h:g4` | `25` |
| `tangentMediumReductionLimit` | `utan.h:g5` | `1e8` |
| `asinTaylorCubic` | `uasncs.h:f1` | `1.66666666666664110590506577996662e-1` |
| `asinTaylorQuintic` | `uasncs.h:f2` | `7.50000000026122686814431784722623e-2` |
| `asinTaylorSeptic` | `uasncs.h:f3` | `4.46428561421059750978517350006940e-2` |
| `asinTaylorNonic` | `uasncs.h:f4` | `3.03821268582119319911193410625235e-2` |
| `asinTaylorUndecic` | `uasncs.h:f5` | `2.23551211026525610742786300334557e-2` |
| `asinTaylorTredecic` | `uasncs.h:f6` | `1.81382903404565056280372531963613e-2` |
| `asinRootRoundBias` | `uasncs.h:t24` | `0x1p24` |
| `acosRootSplitFactor` | `uasncs.h:t27` | `0x1p27` |
| `inverseRootConstant` | `uasncs.h:rt0` | `9.99999999859990725855365213134618e-1` |
| `inverseRootLinear` | `uasncs.h:rt1` | `4.99999999495955425917856814202739e-1` |
| `inverseRootQuadratic` | `uasncs.h:rt2` | `3.75017500867345182581453026130850e-1` |
| `inverseRootCubic` | `uasncs.h:rt3` | `3.12523626554518656309172508769531e-1` |
| `atanTaylorCubic` | `atnat.h:d3` | `-0x1.5555555555555p-2` |
| `atanTaylorQuintic` | `atnat.h:d5` | `0x1.99999999997fdp-3` |
| `atanTaylorSeptic` | `atnat.h:d7` | `-0x1.24924923f7603p-3` |
| `atanTaylorNonic` | `atnat.h:d9` | `0x1.c71c6e5129a3bp-4` |
| `atanTaylorUndecic` | `atnat.h:d11` | `-0x1.7458022b13c25p-4` |
| `atanTaylorTredecic` | `atnat.h:d13` | `0x1.375f08b31cbcep-4` |
| `atanIdentityLimit` | `atnat.h:A/a` | `0x1.bb67a00000000p-27` |
| `atanTaylorLimit` | `atnat.h:B/b` | `0.0625` |
| `atanReciprocalLimit` | `atnat.h:C/c` | `1` |
| `atanReciprocalTableLimit` | `atnat.h:D/d` | `16` |
| `atanSaturationLimit` | `atnat.h:E/e` | `0x1.49ff200000000p52` |
| `atan2ScaleUp` | `atnat2.h:two500` | `0x1p500` |
| `atan2ScaleDown` | `atnat2.h:twom500` | `0x1p-500` |
| `largeReductionInputScale` | `branred.h:tm600` | `0x1p-600` |
| `base24FractionScale` | `branred.h:tm24` | `0x1p-24` |
| `absoluteHighMask` | `s_sin.c:sign_mask` | `0x7fffffff` |
| `exponentHighMask` | `s_sin.c:infinity_high_word` | `0x7ff00000` |
| `fractionHighMask` | `e_asin.c:fraction_mask` | `0x000fffff` |
| `sineIdentityHighWord` | `s_sin.c:identity_threshold` | `(1023 - 26) << 20` |
| `cosineUnityHighWord` | `s_sin.c:cos_unity_threshold` | `(1023 - 27) << 20` |
| `sincosDirectTableHighWord` | `s_sin.c:direct_table_threshold` | `0x3feb6000` |
| `sincosComplementHighWord` | `s_sin.c:complement_threshold` | `0x400368fd` |
| `sincosMediumReductionHighWord` | `s_sin.c:medium_reduction_threshold` | `0x419921fb` |
| `acosUnityHighWord` | `e_asin.c:acos_unity_threshold` | `0x3c880000` |
| `asinTaylorHighWord` | `e_asin.c:asin_Taylor_threshold` | `(1023 - 3) << 20` |
| `asinQuarterHighWord` | `e_asin.c:quarter_threshold` | `(1023 - 2) << 20` |
| `asinHalfHighWord` | `e_asin.c:half_threshold` | `(1023 - 1) << 20` |
| `asinThreeQuarterHighWord` | `e_asin.c:three_quarter_threshold` | `0x3fe80000` |
| `asinFirstUpperHighWord` | `e_asin.c:0.921875_threshold` | `0x3fed8000` |
| `asinSecondUpperHighWord` | `e_asin.c:0.953125_threshold` | `0x3fee8000` |
| `asinTableUpperHighWord` | `e_asin.c:0.96875_threshold` | `0x3fef0000` |
| `domainUnitHighWord` | `e_asin.c:unit_threshold` | `1023 << 20` |
| `atan2ExponentGapLimit` | `e_atan2.c:ep/em` | `57 << 20` |
| `asinQuarterSegmentOffset` | `e_asin.c:352` | `352` |
| `asinHalfSegmentOffset` | `e_asin.c:1056` | `1056` |
| `asinThreeQuarterSegmentOffset` | `e_asin.c:992` | `992` |
| `asinFirstUpperSegmentOffset` | `e_asin.c:884` | `884` |
| `asinSecondUpperSegmentOffset` | `e_asin.c:768` | `768` |
| `atanTableFirstIndex` | `s_atan.c:16` | `16` |
| `largeReductionScaleHighWord` | `branred.h:t576_high_word` | `0x63f00000` |
| `largeReductionExponentBias` | `branred.c:450` | `450` |
| `base24Digits` | `branred.c:24` | `24` |
| `quietNaNBit` | `IEEE754:quiet_bit` | `1 << 51` |
| `positiveQuietNaNBits` | `glibc:asin/acos_domain_NaN` | `0x7ff8000000000000` |
| `negativeQuietNaNBits` | `glibc:sin/cos/tan_domain_NaN` | `0xfff8000000000000` |

## FMA assembly audit

The audit used GCC 13.3.0 (Ubuntu 13.3.0-6ubuntu2~24.04.1) and GNU objdump 2.42 on x86-64. The source hashes were checked against the manifest. The numerical C files were compiled separately and inspected with `objdump -d -Mintel -S`, including the final RTL-generated instructions. GIMPLE alone was not used as the contraction authority.

The [x86-64 multiarch Makefile](https://github.com/bminor/glibc/blob/glibc-2.39/sysdeps/x86_64/fpu/multiarch/Makefile) adds `-mfma -mavx2` to e_asin-fma, e_atan2-fma, s_atan-fma, s_sin-fma, and s_tan-fma. These wrappers include the audited dbl-64 bodies. The body audit used `-O2 -mfma -mavx2 -ffp-contract=fast -frounding-math -g`.

Branred is a deliberate exception. The [dbl-64 Makefile](https://github.com/bminor/glibc/blob/glibc-2.39/sysdeps/ieee754/dbl-64/Makefile) says `CFLAGS-branred.c += $(config-cflags-nofma)`. The sin/tan FMA wrappers call the ordinary __branred, not an FMA copy. It was audited with `-O2 -ffp-contract=off -frounding-math -mprefer-vector-width=128 -g`; its object contains no FMA instruction. Compiling this helper with the caller's FMA flags would change its split operations and produce incorrect large-argument results.

The stand-in headers omit floating-point environment save/restore and underflow forcing. The objects retain unresolved errno calls; the audit disassembles the objects without linking or executing those calls. No polynomial, table, reduction, or double-length operation is replaced. The numerical instruction boundaries below refer to the compiled dbl-64 function symbols, before wrapper renaming. Counts are static occurrences in the object, including compiler-duplicated paths; they are not per-call operation counts.

| C body | Scalar FMA instructions | Audited regions and Go implementation |
| --- | ---: | --- |
| s_sin.c: __sin | 70 | Taylor Horner chain, two table polynomials and correction sums, both medium-reduction correction products, and all quadrant paths; taylorSin/doSin/doCos/reduceSinCos. |
| s_sin.c: __cos | 70 | The same polynomial/reduction families plus the complement-range double-length correction; Cos/doSinCos and the shared helpers. |
| s_tan.c: __tan | 45 | Taylor/table chains, short and medium argument reduction, and the EMULV/DIV2 quotient residual; Tan/tanPolynomial/tanTableResult/dividePair. |
| e_asin.c: __ieee754_asin | 61 | All five table segments, both small-argument paths, reciprocal-root refinement, denominator, and final half-pi correction; asinApprox/asinRoot/Asin. |
| e_asin.c: __ieee754_acos | 60 | The same segment/refinement families with the two-FMA root split and the separately rounded branch-shared final product; acosRoot/Acos. |
| s_atan.c: __atan | 28 | Both Horner chains, table correction, exact reciprocal residual and its two uses; atanPolynomial/atanTablePolynomial/atanReciprocal/Atan. |
| e_atan2.c: __ieee754_atan2 | 47 | Both ratio residuals, direct Taylor/table corrections, and the three quadrant corrections; atan2Ratio/atan2Direct/atan2AddOffset/atan2SubtractOffset. |
| branred.c: __branred | 0 | Both mantissa splits, base-24 chunks, high/low products and final compensation stay separate; reducePart/reduceLarge. |

The following instruction sites distinguish contractions that cross assignments or reuse one product:

| Symbol and object offset | ASM evidence | Explicit Go operation |
| --- | --- | --- |
| __tan+0xd3, +0xdc | Two vfnmadd instructions consume xn and mp3, without an intervening rounded xn*mp3. | `a = math.FMA(-xn, piHalfCorrection, t1)`; `da = math.FMA(-xn, piHalfCorrection, t1-a)`. |
| __tan+0x249, +0x252, +0x25b, +0x264 | Each use of xn*pp3 and xn*pp4 contracts separately. | The medium-reduction t/da/a/da chain has four math.FMA calls. |
| __sin+0x196, +0x19f, +0x1ac, +0x1b5 | Both uses of each medium-reduction correction product contract after RTL. | reduceSinCos uses four math.FMA calls for the two compensation stages. |
| __ieee754_asin+0x437 and __ieee754_acos+0x7fe | The refined-root product contracts with the denominator's +y. | `math.FMA(factor, c, y)` in asinRoot/acosRoot. |
| __ieee754_acos+0x7f0, +0x7f5 | Both appearances of t27*c contract, despite the repeated product. | Two nested math.FMA calls in acosRoot. |
| __ieee754_acos+0x849 | The final `(y+cc)*p` is a vmulsd shared across the sign branch. | Acos keeps `product := float64(p*(y+cc))` before either correction. |
| __atan+0x170 and +0x24c | EMULV retains a rounded multiply and uses vfmsub for its residual. | atanReciprocal keeps `float64(w*u)` and `math.FMA(w,u,-t1)`. |
| __atan+0x27c, +0x281 | Both reciprocal/tail corrections contract with the low half-pi term. | The two math.FMA calls in the final Atan reciprocal branch. |
| __ieee754_atan2+0x555, +0x5ff, +0x97b | Offset Taylor products remain vmulsd across the exact-add/subtract control flow. | Add/subtract offsets keep both `float64(u*v)` and the separately rounded polynomial product. |

Every floating-point multiplication in the translated Go code is either an explicit math.FMA operand or wrapped as `float64(a*b)`. Sums and differences keep their original parenthesization without redundant float64 conversions. The resulting port passes all 986775 recorded Redis comparisons on Windows, including the expanded pole neighborhoods; final cloud and CI evidence is reported in the PR.

The audit objects, disassembly, source/build hashes, compiler versions and reproducible build command are saved in `work/asm-audit-20261024.tar.gz` for the review record.
