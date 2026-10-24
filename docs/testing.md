# Testing

`TestRedisGrid` checks 986775 inputs against recorded Redis 7.2.16 answers. Every supported OS executes the same compressed fixture and deterministic generator. The fixture was captured on the Azure reference machine through Redis EVAL. All 700000 original answers remain unchanged; the additional 286775 answers extend the boundary coverage.

| Function | Recorded answers |
| --- | ---: |
| sin | 181930 |
| cos | 181930 |
| tan | 222895 |
| asin | 100010 |
| acos | 100010 |
| atan | 100000 |
| atan2 | 100000 |

The generator in `internal/grid` describes each function by its name, evaluator, generator, and count. Separate input families cover ordinary circular functions, inverse domains, and two-argument atan. Its SplitMix64 inputs have no package state. The original inputs cover ordinary values, table intervals, range boundaries and their neighbors, signed zero, subnormals, infinities, NaNs, and large trigonometric arguments beyond 1e300. `Atan2` includes all quadrants, scale changes, and large exponent differences.

For sin and cos, extra grids contain five consecutive binary64 values centered on every multiple of pi/2 and pi for integer multipliers from -4096 through 4096. Tan additionally includes multiples of pi/4. These neighborhoods include both neighbors up to two ULP away. Asin and acos add the corresponding five-value neighborhoods of -1 and +1, including values outside the domain. Every floating product used to generate inputs has an explicit rounding barrier, so input bits agree across architectures.

`internal/cmd/capture` requires Redis 7.2.16. It evaluates batches of 500 inputs and records Redis `string.format('%.17g', ...)` replies as named TSV rows. `-verify` checks every line against a fresh Redis execution and rejects missing, reordered, or extra records. To regenerate the fixture, omit `-verify`; capture replaces the fixture only after every response succeeds. Expected answers must never be edited by hand.

CI has the required `test (ubuntu-latest)`, `test (windows-latest)`, `test (macos-latest)`, `lint`, `race`, and `reference` jobs. The macOS job runs on arm64. The root package coverage threshold is 85%. The module dependency check requires `go list -m all` to contain one line.

The integration test checks both `testdata/snippets.json` and `testdata/libm.json` from the Lua checkout pinned in `.github/workflows/ci.yml`. This includes the 20 reviewer L1 snippets. It calls the public `github.com/acneism/lua/luatest.RunCorpus(t, files, options)` helper with all seven functions supplied through `Options.Math`. The helper uses Lua's original wrapper, environment, limits, and normalization rules. The integration command creates a separate temporary test module with local replacements; it never patches Lua sources or corpus files. Neither production module gains a dependency or replace directive. To run it against a checkout containing the A3 helper:

```text
go run ./internal/cmd/integration -lua ../lua
```

`FuzzTrig` evaluates every function twice on arbitrary binary64 bit patterns and requires identical bits without a panic. PR CI runs it for 30 seconds; the nightly workflow runs it for ten minutes.

The Go mutation command creates isolated copies of the sources and fixture. Each mutation must fail `TestRedisGrid`: remove the atan table correction FMA, change the large-argument quadrant in `branred`, flip one bit in a sine table entry, or revert both tangent small-reduction correction FMAs. It leaves production sources untouched and saves the failing test logs in `work/mutations/`.

```text
go run ./internal/cmd/mutate
```
