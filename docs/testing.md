# Testing

`TestRedisGrid` checks 100000 inputs per function against recorded Redis 7.2.16 answers, for 700000 comparisons. Every supported OS executes the same compressed fixture and generator. The fixture was captured on the Azure reference machine through Redis EVAL, not through the port or native libm calls.

The deterministic generator is `internal/grid`. Its SplitMix64 inputs have no package state. It covers ordinary values, table intervals, range boundaries and their neighbors, signed zero, subnormals, infinities, NaNs, and large trigonometric arguments beyond 1e300. `Asin` and `Acos` concentrate samples near both ends of their domains. `Atan2` includes all quadrants, scale changes, and large exponent differences.

`internal/cmd/capture` requires Redis 7.2.16. It evaluates batches of 500 inputs and records Redis `string.format('%.17g', ...)` replies. `-verify` checks every line against a fresh Redis execution and rejects missing, reordered, or extra records. To regenerate the fixture, omit `-verify`. Expected answers must never be edited by hand.

CI has the required `test (ubuntu-latest)`, `test (windows-latest)`, `test (macos-latest)`, `lint`, `race`, and `reference` jobs. The macOS job runs on arm64. The root package coverage threshold is 85%. The module dependency check requires `go list -m all` to contain one line.

The integration job pins `lua` to commit `3ab696bab0ed4b85f5de1fe2ca09ad1df301d18a`. `scripts/integration.py` copies that checkout into a temporary directory, injects all seven functions into its existing `TestCorpus`, and appends `testdata/libm.json` to the temporary corpus. Original answers remain unchanged. The temporary Lua module alone receives a local test dependency on this module; neither production module gains a dependency or replace directive. Each OS executes all 4129 Lua corpus snippets and the 16 precise libm snippets.

`FuzzTrig` evaluates every function twice on arbitrary binary64 bit patterns and requires identical bits without a panic. PR CI runs it for 30 seconds; the nightly workflow runs it for ten minutes.

`scripts/mutate.py` creates temporary copies and checks three mutations. Removing the FMA from the atan table correction, changing the large-argument quadrant in `branred`, and flipping one bit of a sine table constant must each fail `TestRedisGrid`. The script leaves production sources untouched and saves logs in `work/`.
