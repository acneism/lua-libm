package libm_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	libm "github.com/acneism/lua-libm"
	"github.com/acneism/lua-libm/internal/grid"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func evaluate(function int, x, y float64) float64 {
	if function == 6 {
		return libm.Atan2(x, y)
	}
	return [6]func(float64) float64{libm.Sin, libm.Cos, libm.Tan, libm.Asin, libm.Acos, libm.Atan}[function](x)
}

func resultText(x float64) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return grid.Argument(x)
	}
	return fmt.Sprintf("%.17g", x)
}

func TestRedisGrid(t *testing.T) {
	file, err := os.Open("testdata/grid.tsv.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()
	scan := bufio.NewScanner(reader)
	var count [grid.Functions]int
	var mismatches [grid.Functions]int
	for scan.Scan() {
		fields := strings.Split(scan.Text(), "\t")
		if len(fields) != 3 {
			t.Fatal("invalid fixture row")
		}
		function, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		index, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		x, y := grid.Input(function, index)
		got := resultText(evaluate(function, x, y))
		want := fields[2]
		count[function]++
		if got != want {
			if mismatches[function] < 8 {
				t.Errorf("%s[%d](%.17g,%.17g): got %s want %s", grid.Name(function), index, x, y, got, want)
			}
			mismatches[function]++
		}
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	for function, n := range count {
		if n != grid.Count {
			t.Errorf("%s: got %d rows want %d", grid.Name(function), n, grid.Count)
		}
		t.Logf("%s: %d inputs, %d differences", grid.Name(function), n, mismatches[function])
	}
}

func FuzzTrig(f *testing.F) {
	for _, bits := range []uint64{0, 1, 1 << 63, 0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000000, 0x3ff0000000000000, 0x3fefffffffffffff, 0x7fefffffffffffff, 0x7e37e43c8800759c} {
		f.Add(bits, bits^0x8000000000000000)
	}
	f.Fuzz(func(t *testing.T, xBits, yBits uint64) {
		x, y := math.Float64frombits(xBits), math.Float64frombits(yBits)
		for function := 0; function < grid.Functions; function++ {
			a := math.Float64bits(evaluate(function, x, y))
			b := math.Float64bits(evaluate(function, x, y))
			if a != b {
				t.Fatalf("%s(%x,%x) changed from %x to %x", grid.Name(function), xBits, yBits, a, b)
			}
		}
	})
}
