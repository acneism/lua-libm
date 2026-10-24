package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type mutation struct {
	name, file string
	rewrite    func(*ast.File) error
}

func main() {
	root := flag.String("root", ".", "repository to copy into isolated mutation workspaces")
	output := flag.String("output", "work/mutations", "directory for mutation logs")
	flag.Parse()
	if err := run(*root, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root, output string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(output, 0755); err != nil {
		return err
	}
	mutations := []mutation{
		{"fma", "atan.go", removeAtanTableFMA},
		{"branred", "reduce.go", changeQuadrant},
		{"table", "tables.go", changeTable},
		{"tan-poles", "tan.go", revertTangentCorrection},
	}
	for _, item := range mutations {
		if err := checkMutation(root, output, item); err != nil {
			return err
		}
		fmt.Printf("%s -> TestRedisGrid: caught\n", item.name)
	}
	return nil
}

func checkMutation(root, output string, item mutation) error {
	temporary, err := os.MkdirTemp(output, "mutation-")
	if err != nil {
		return err
	}
	defer func() {
		relative, relErr := filepath.Rel(output, temporary)
		if relErr == nil && filepath.IsLocal(relative) && strings.HasPrefix(filepath.Base(temporary), "mutation-") {
			_ = os.RemoveAll(temporary)
		}
	}()
	for _, directory := range []string{".", "internal", "testdata"} {
		if err := copySources(root, temporary, directory); err != nil {
			return err
		}
	}
	path := filepath.Join(temporary, item.file)
	set := token.NewFileSet()
	source, err := parser.ParseFile(set, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if err = item.rewrite(source); err != nil {
		return fmt.Errorf("%s: %w", item.name, err)
	}
	var formatted bytes.Buffer
	if err = format.Node(&formatted, set, source); err != nil {
		return err
	}
	if err = os.WriteFile(path, formatted.Bytes(), 0644); err != nil {
		return err
	}
	command := exec.Command("go", "test", ".", "-run", "^TestRedisGrid$", "-count=1")
	command.Dir = temporary
	result, runErr := command.CombinedOutput()
	if err = os.WriteFile(filepath.Join(output, item.name+".log"), result, 0644); err != nil {
		return err
	}
	if runErr == nil || !bytes.Contains(result, []byte("--- FAIL: TestRedisGrid")) {
		return fmt.Errorf("%s survived or failed outside TestRedisGrid: %v\n%s", item.name, runErr, result)
	}
	return nil
}

func copySources(root, destination, directory string) error {
	return filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if directory == "." && entry.IsDir() && relative != "." {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return os.MkdirAll(filepath.Join(destination, relative), 0755)
		}
		if directory == "." && filepath.Ext(path) != ".go" && entry.Name() != "go.mod" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(destination, relative), data, 0644)
	})
}

func isFMA(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	name, ok := selector.X.(*ast.Ident)
	return ok && name.Name == "math" && selector.Sel.Name == "FMA" && len(call.Args) == 3
}

func removeAtanTableFMA(file *ast.File) error {
	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if changed || !ok || !isFMA(call) {
			return !changed
		}
		polynomial, ok := call.Args[1].(*ast.CallExpr)
		if !ok {
			return true
		}
		name, ok := polynomial.Fun.(*ast.Ident)
		if !ok || name.Name != "atanTablePolynomial" {
			return true
		}
		product := &ast.CallExpr{Fun: ast.NewIdent("float64"), Args: []ast.Expr{&ast.BinaryExpr{X: call.Args[0], Op: token.MUL, Y: call.Args[1]}}}
		*call = ast.CallExpr{Fun: ast.NewIdent("float64"), Args: []ast.Expr{&ast.BinaryExpr{X: product, Op: token.ADD, Y: call.Args[2]}}}
		changed = true
		return false
	})
	if !changed {
		return fmt.Errorf("FMA call missing")
	}
	return nil
}

func changeQuadrant(file *ast.File) error {
	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		assignment, ok := node.(*ast.AssignStmt)
		if changed || !ok || len(assignment.Rhs) != 1 {
			return !changed
		}
		expression, ok := assignment.Rhs[0].(*ast.BinaryExpr)
		if !ok || expression.Op != token.AND {
			return true
		}
		constant, ok := expression.Y.(*ast.BasicLit)
		if !ok || constant.Value != "3" {
			return true
		}
		call, ok := expression.X.(*ast.CallExpr)
		if !ok {
			return true
		}
		name, ok := call.Fun.(*ast.Ident)
		if !ok || name.Name != "int" {
			return true
		}
		expression.X = &ast.ParenExpr{X: &ast.BinaryExpr{X: call, Op: token.ADD, Y: &ast.BasicLit{Kind: token.INT, Value: "1"}}}
		changed = true
		return false
	})
	if !changed {
		return fmt.Errorf("quadrant assignment missing")
	}
	return nil
}

func changeTable(file *ast.File) error {
	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		value, ok := node.(*ast.ValueSpec)
		if !ok || len(value.Names) != 1 || value.Names[0].Name != "sincosTableBits" {
			return true
		}
		table, ok := value.Values[0].(*ast.CompositeLit)
		if !ok || len(table.Elts) <= 32 {
			return false
		}
		entry, ok := table.Elts[32].(*ast.CompositeLit)
		if !ok || len(entry.Elts) == 0 {
			return false
		}
		field := entry.Elts[0]
		if keyed, ok := field.(*ast.KeyValueExpr); ok {
			field = keyed.Value
		}
		literal, ok := field.(*ast.BasicLit)
		if !ok {
			return false
		}
		bits, err := strconv.ParseUint(literal.Value, 0, 64)
		if err != nil {
			return false
		}
		literal.Value = fmt.Sprintf("0x%016x", bits^1)
		changed = true
		return false
	})
	if !changed {
		return fmt.Errorf("sine table entry missing")
	}
	return nil
}

func revertTangentCorrection(file *ast.File) error {
	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || !isFMA(call) {
			return true
		}
		coefficient, ok := call.Args[1].(*ast.Ident)
		if !ok || coefficient.Name != "piHalfCorrection" {
			return true
		}
		product := &ast.CallExpr{Fun: ast.NewIdent("float64"), Args: []ast.Expr{&ast.BinaryExpr{X: call.Args[0], Op: token.MUL, Y: call.Args[1]}}}
		*call = ast.CallExpr{Fun: ast.NewIdent("float64"), Args: []ast.Expr{&ast.BinaryExpr{X: product, Op: token.ADD, Y: call.Args[2]}}}
		changed = true
		return false
	})
	if !changed {
		return fmt.Errorf("tangent correction FMA missing")
	}
	return nil
}
