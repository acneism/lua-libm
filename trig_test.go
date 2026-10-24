package libm_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/acneism/lua-libm/internal/grid"
)

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
	functions := grid.Functions()
	byName := make(map[string]grid.Function, len(functions))
	for _, function := range functions {
		byName[function.Name] = function
	}
	counts := make(map[string]int, len(functions))
	mismatches := make(map[string]int, len(functions))
	scan := bufio.NewScanner(reader)
	for scan.Scan() {
		fields := strings.Split(scan.Text(), "\t")
		if len(fields) != 3 {
			t.Fatal("invalid fixture row")
		}
		function, ok := byName[fields[0]]
		if !ok {
			t.Fatalf("unknown fixture function %q", fields[0])
		}
		index, err := strconv.Atoi(fields[1])
		if err != nil || index != counts[function.Name] || index >= function.Count {
			t.Fatalf("invalid %s index %q", function.Name, fields[1])
		}
		x, y := function.Generate(index)
		got := resultText(function.Evaluate(x, y))
		want := fields[2]
		counts[function.Name]++
		if got != want {
			if mismatches[function.Name] < 8 {
				t.Errorf("%s[%d](%.17g,%.17g): got %s want %s", function.Name, index, x, y, got, want)
			}
			mismatches[function.Name]++
		}
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	for _, function := range functions {
		n := counts[function.Name]
		if n != function.Count {
			t.Errorf("%s: got %d rows want %d", function.Name, n, function.Count)
		}
		t.Logf("%s: %d inputs, %d differences", function.Name, n, mismatches[function.Name])
	}
}

func FuzzTrig(f *testing.F) {
	for _, bits := range []uint64{0, 1, 1 << 63, 0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000000, 0x3ff0000000000000, 0x3fefffffffffffff, 0x7fefffffffffffff, 0x7e37e43c8800759c} {
		f.Add(bits, bits^0x8000000000000000)
	}
	functions := grid.Functions()
	f.Fuzz(func(t *testing.T, xBits, yBits uint64) {
		x, y := math.Float64frombits(xBits), math.Float64frombits(yBits)
		for _, function := range functions {
			a := math.Float64bits(function.Evaluate(x, y))
			b := math.Float64bits(function.Evaluate(x, y))
			if a != b {
				t.Fatalf("%s(%x,%x) changed from %x to %x", function.Name, xBits, yBits, a, b)
			}
		}
	})
}
