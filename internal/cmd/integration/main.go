package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func main() {
	lua := flag.String("lua", "", "Lua checkout containing the public luatest.RunCorpus helper")
	root := flag.String("root", ".", "lua-libm checkout to test through Options.Math")
	flag.Parse()
	if err := run(*lua, *root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(lua, root string) error {
	if lua == "" {
		return fmt.Errorf("-lua must name a Lua checkout")
	}
	lua, err := filepath.Abs(lua)
	if err != nil {
		return err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	files := []string{filepath.Join(lua, "testdata", "snippets.json"), filepath.Join(lua, "testdata", "libm.json")}
	counts := make([]int, len(files))
	for index, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var entries []json.RawMessage
		if err = json.Unmarshal(data, &entries); err != nil {
			return err
		}
		counts[index] = len(entries)
	}
	temporary, err := os.MkdirTemp("", "lua-libm-integration-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(temporary) }()
	module := fmt.Sprintf(`module integration

go 1.26.1

require (
	github.com/acneism/lua v0.0.0
	github.com/acneism/lua-libm v0.0.0
)

replace github.com/acneism/lua => %s
replace github.com/acneism/lua-libm => %s
`, strconv.Quote(filepath.ToSlash(lua)), strconv.Quote(filepath.ToSlash(root)))
	if err = os.WriteFile(filepath.Join(temporary, "go.mod"), []byte(module), 0644); err != nil {
		return err
	}
	test := fmt.Sprintf(`package integration_test

import (
	"testing"

	"github.com/acneism/lua"
	libm "github.com/acneism/lua-libm"
	"github.com/acneism/lua/luatest"
)

func TestCorpus(t *testing.T) {
	luatest.RunCorpus(t, []string{%s, %s}, lua.Options{Math: &lua.Math{
		Sin: libm.Sin, Cos: libm.Cos, Tan: libm.Tan,
		Asin: libm.Asin, Acos: libm.Acos, Atan: libm.Atan, Atan2: libm.Atan2,
	}})
}
`, strconv.Quote(files[0]), strconv.Quote(files[1]))
	if err = os.WriteFile(filepath.Join(temporary, "integration_test.go"), []byte(test), 0644); err != nil {
		return err
	}
	command := exec.Command("go", "test", ".", "-run", "^TestCorpus$", "-count=1")
	command.Dir = temporary
	command.Env = append(os.Environ(), "GOWORK=off")
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	if err = command.Run(); err != nil {
		return err
	}
	fmt.Printf("integration: %d corpus snippets and %d exact libm snippets through luatest.RunCorpus\n", counts[0], counts[1])
	return nil
}
