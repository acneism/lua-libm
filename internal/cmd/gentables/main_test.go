package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseTable(t *testing.T) {
	text := "#ifdef BIG_ENDI\nstatic const union {int i[4];double x[2];} table = {.i={0x3ff00000,0x00000000,0xbff00000,0x00000000}};\n#else\ninvalid\n#endif"
	got, err := parseTable(text, "table", false)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []uint64{0x3ff0000000000000, 0xbff0000000000000}) {
		t.Fatalf("got %x", got)
	}
	got, err = parseTable("static const double inroot[] = { 1.0, 0.5, -2.0 };", "inroot", true)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []uint64{0x3ff0000000000000, 0x3fe0000000000000, 0xc000000000000000}) {
		t.Fatalf("got %x", got)
	}
	for _, invalid := range []string{"missing", "table=", "table={0x1};", "#ifdef BIG_ENDI\ntable={0x1};"} {
		if _, err := parseTable(invalid, "table", false); err == nil {
			t.Fatalf("accepted %q", invalid)
		}
	}
}

func TestCheckedSources(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "source.c")
	data := []byte("source")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(directory, "sources.sha256")
	if err := os.WriteFile(manifest, fmt.Appendf(nil, "%x  source.c\n", sha256.Sum256(data)), 0600); err != nil {
		t.Fatal(err)
	}
	files, err := checkedSources(directory, manifest)
	if err != nil || files["source.c"] != string(data) {
		t.Fatalf("files=%v err=%v", files, err)
	}
	if err := os.WriteFile(path, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := checkedSources(directory, manifest); err == nil {
		t.Fatal("accepted a modified source")
	}
}
