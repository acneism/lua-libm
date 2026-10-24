# Porting

The source version is [glibc 2.39](https://sourceware.org/git/?p=glibc.git;a=tree;f=sysdeps/ieee754/dbl-64;hb=glibc-2.39). The downloaded files and their SHA-256 digests are recorded in [sources.sha256](sources.sha256). The tag is also available in the [glibc source mirror](https://github.com/bminor/glibc/tree/glibc-2.39/sysdeps/ieee754/dbl-64).

| Go source | glibc source |
| --- | --- |
| sin.go | s_sin.c, usncs.h |
| reduce.go | branred.c, branred.h |
| tan.go | s_tan.c, utan.h |
| asin.go | e_asin.c, uasncs.h |
| atan.go | s_atan.c, atnat.h |
| atan2.go | e_atan2.c, atnat2.h |
| arithmetic.go | dla.h, mydefs.h |
| tables.go | sincostab.c, utan.tbl, asincos.tbl, root.tbl, powtwo.tbl, uatan.tbl, branred.h |

C unions become `math.Float64bits` and `math.Float64frombits`. Each table stores the original binary64 bit patterns in a constant byte string. `binary.LittleEndian.Uint64` reads a value without a mutable global array or initialization step. Decimal coefficients are rounded to binary64 typed constants.

Repeated polynomial evaluation and tangent reduction blocks become helper functions. The table indexes, Horner order, double-length correction formulas, range thresholds, and quadrant selection remain the same.

FMA placement follows the glibc FMA variant compiled by GCC with `-mfma -ffp-contract=fast`. The translation was checked against GCC's optimized GIMPLE. Contraction can cross C statement boundaries: the table correction in `asin`, the reciprocal correction in `atan`, and the refined-root denominator are examples. A product shared across branches can prevent contraction, as in the final `acos` correction. Explicit `float64` conversions force every operation that glibc leaves separate.

NaN selection follows the public glibc entry points used by Redis, including their domain handling. The Go API does not reproduce floating-point environment or errno updates. README and NOTICE identify these translation changes and preserve the LGPL notices outside Go source comments.
