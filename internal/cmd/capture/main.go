package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/acneism/lua-libm/internal/grid"
	"github.com/acneism/lua-libm/internal/redisref"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (resultErr error) {
	address := flag.String("addr", os.Getenv("REDIS_REF_ADDR"), "Redis 7.2.16 address as host:port; defaults to REDIS_REF_ADDR")
	output := flag.String("output", "testdata/grid.tsv.gz", "gzip TSV fixture to capture or verify")
	verify := flag.Bool("verify", false, "compare every recorded answer with Redis without writing the fixture")
	flag.Parse()
	if *address == "" {
		return fmt.Errorf("set REDIS_REF_ADDR or -addr")
	}
	client, err := redisref.Dial(*address)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()
	info, err := client.Command("INFO", "server")
	if err != nil {
		return err
	}
	server, ok := info.(string)
	if !ok || !strings.Contains(server, "redis_version:7.2.16\r\n") {
		return fmt.Errorf("redis 7.2.16 required")
	}
	var destination io.Writer
	var expected *bufio.Scanner
	var captured *os.File
	var compressed *gzip.Writer
	if *verify {
		file, err := os.Open(*output)
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); resultErr == nil {
				resultErr = err
			}
		}()
		reader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer func() { _ = reader.Close() }()
		expected = bufio.NewScanner(reader)
	} else {
		file, err := os.CreateTemp(filepath.Dir(*output), ".capture-*")
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
			_ = os.Remove(file.Name())
		}()
		writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
		if err != nil {
			return err
		}
		defer func() {
			if err := writer.Close(); resultErr == nil {
				resultErr = err
			}
		}()
		destination = writer
		captured, compressed = file, writer
	}
	for _, function := range grid.Functions() {
		name := function.Name
		script := "local r={} for i=1,#ARGV,2 do r[#r+1]=string.format('%.17g',math." + name + "(tonumber(ARGV[i]),tonumber(ARGV[i+1]))) end return r"
		for base := 0; base < function.Count; base += 500 {
			n := min(500, function.Count-base)
			args := []string{"EVAL", script, "0"}
			for i := 0; i < n; i++ {
				x, y := function.Generate(base + i)
				args = append(args, grid.Argument(x), grid.Argument(y))
			}
			response, err := client.Command(args...)
			if err != nil {
				return err
			}
			rows, ok := response.([]any)
			if !ok || len(rows) != n {
				return fmt.Errorf("invalid reference response for %s", name)
			}
			for i, row := range rows {
				want, ok := row.(string)
				if !ok {
					return fmt.Errorf("non-string reference response")
				}
				line := fmt.Sprintf("%s\t%d\t%s", name, base+i, want)
				if *verify {
					if !expected.Scan() || expected.Text() != line {
						return fmt.Errorf("reference differs: %s %d: %q != %q", name, base+i, expected.Text(), line)
					}
				} else {
					if _, err := fmt.Fprintln(destination, line); err != nil {
						return err
					}
				}
			}
		}
		fmt.Printf("%s: %d Redis answers\n", name, function.Count)
	}
	if *verify {
		if expected.Scan() {
			return fmt.Errorf("extra fixture rows")
		}
		return expected.Err()
	}
	if err := compressed.Close(); err != nil {
		return err
	}
	if err := captured.Close(); err != nil {
		return err
	}
	return os.Rename(captured.Name(), *output)
}
