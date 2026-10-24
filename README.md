# libm

This module ports the IBM Accurate Mathematical Library routines from glibc 2.39 to Go. It provides `Sin`, `Cos`, `Tan`, `Asin`, `Acos`, `Atan`, and `Atan2` for Redis 7.2.16 compatible Lua execution.

The target is glibc on x86-64 with FMA enabled. Every function uses the same algorithm and explicit `math.FMA` operations. Explicit `float64` conversions preserve the remaining rounding steps on other architectures. The module uses only the Go standard library. Its package variables are read-only lookup tables; there is no shared mutable state, cgo, assembly, `unsafe`, or `reflect` in production code.

```go
package main

import (
	"github.com/acneism/lua"
	libm "github.com/acneism/lua-libm"
)

func main() {
	state := lua.NewState(lua.Options{Math: &lua.Math{
		Sin: libm.Sin, Cos: libm.Cos, Tan: libm.Tan,
		Asin: libm.Asin, Acos: libm.Acos, Atan: libm.Atan,
		Atan2: libm.Atan2,
	}})
	_ = state
}
```

`Atan2(y, x)` takes arguments in that order. Each other function takes one `float64`. Functions have no state and can run concurrently.

Run `go test ./...` and `go vet ./...`. The checked-in fixture contains 986775 Redis answers, including 100000 original inputs per function and neighborhoods of trigonometric reduction and domain boundaries. To verify every recorded answer against a running Redis 7.2.16, run:

```text
go run ./internal/cmd/capture -addr 127.0.0.1:6379 -verify
```

[Testing](docs/testing.md) describes the grids, integration test, fuzz target, and mutation checks. [Porting](docs/porting.md) lists the upstream sources and the changes made during translation. [Compatibility](docs/compatibility.md) describes the numerical contract.

This module is LGPL-2.1-or-later. [LICENSE](LICENSE) contains LGPL 2.1, and [NOTICE](NOTICE) preserves the upstream notices. A host that distributes this module must provide the license text and a link to its source in the host's NOTICE and distribution image. The Apache-2.0 `lua` module does not import it; the host selects it through `Options.Math`.
