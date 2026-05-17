package lint

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTestFiles_SamePackageAsProduction rejects external test packages (`package foo_test`)
// so tests can access unexported helpers and the package stays one compilation unit.
func TestTestFiles_SamePackageAsProduction(t *testing.T) {
	root := RepoRoot(t)
	fset := token.NewFileSet()

	// Per-dir: production package name (any non-_test.go file's package) and the list of
	// test files with their declared package name.
	type testEntry struct {
		path string
		pkg  string
	}
	prodPkg := map[string]string{}
	testFiles := map[string][]testEntry{}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		dir := filepath.Dir(path)
		pkg := f.Name.Name
		if strings.HasSuffix(path, "_test.go") {
			testFiles[dir] = append(testFiles[dir], testEntry{path: path, pkg: pkg})
		} else {
			prodPkg[dir] = pkg
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	for dir, entries := range testFiles {
		prod, ok := prodPkg[dir]
		if !ok {
			// Directory has only test files — no production package to match against.
			continue
		}
		for _, e := range entries {
			if e.pkg != prod {
				rel, _ := filepath.Rel(root, e.path)
				t.Errorf("%s: declares package %q, expected %q (matching production code in same directory). "+
					"External test packages are not allowed — use same-package tests so they can access "+
					"unexported helpers directly.", rel, e.pkg, prod)
			}
		}
	}
}
