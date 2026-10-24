# Compatibility

The target is Redis 7.2.16 running on x86-64 with FMA and glibc. The same Go code runs on amd64 and arm64. It does not select the local system's libm.

Finite results must equal the target binary64 value. The reference test compares `%.17g`, which distinguishes every finite binary64 value, including signed zero. NaNs are checked by class and sign; their payload is not represented by the Redis string. No ULP tolerance or result normalization is used in this module's grids.

The APIs match public libm behavior used by Redis. `Asin` and `Acos` return a positive NaN for arguments outside their domains. NaN inputs preserve the selected argument's sign. `Atan2` selects the second argument's NaN first, as glibc does. Infinite inputs to `Sin`, `Cos`, and `Tan` return a negative NaN. These choices are explicit and independent of the instruction set.

The port assumes IEEE 754 binary64 round to nearest, ties to even. It does not expose C `errno`, floating-point exception flags, or directed rounding modes. These are outside the Lua math API.

The module uses glibc's algorithms rather than a correctly rounded mathematical approximation. Improving a glibc rounding error would change a Redis result and fail the reference contract.
