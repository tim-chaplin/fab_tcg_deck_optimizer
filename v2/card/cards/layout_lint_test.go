package cards_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Tests that NotImplemented and Unplayable markers appear only in their dedicated subpackages.
func TestLayout_MarkersStayInSubpackages(t *testing.T) {
	const (
		notImplDir  = "notimplemented"
		unplayDir   = "unplayable"
		notImplName = "NotImplemented"
		unplayName  = "Unplayable"
	)
	roots := map[string]string{
		notImplName: notImplDir,
		unplayName:  unplayDir,
	}

	// Walk the cards tree and collect, per file, which marker methods it declares (if any).
	type finding struct {
		path    string
		dir     string // "" for top-level cards/, else subdir name
		markers []string
	}
	var findings []finding
	err := filepath.WalkDir("./", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		var markers []string
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Recv.NumFields() == 0 {
				continue
			}
			name := fn.Name.Name
			if name == notImplName || name == unplayName {
				markers = append(markers, name)
			}
		}
		if markers == nil {
			return nil
		}
		// Determine which subdir (or top-level) this file lives in.
		dir := ""
		if rel, err := filepath.Rel("./", path); err == nil {
			parts := strings.Split(rel, string(filepath.Separator))
			if len(parts) > 1 {
				dir = parts[0]
			}
		}
		findings = append(findings, finding{path: path, dir: dir, markers: markers})
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	for _, f := range findings {
		for _, m := range f.markers {
			wantDir := roots[m]
			if f.dir != wantDir {
				t.Errorf("%s: declares %s() but lives outside %s/ — move the file there.",
					f.path, m, wantDir)
			}
		}
	}
}
