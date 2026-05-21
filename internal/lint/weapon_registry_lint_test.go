package lint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Tests that every implemented weapon is in registry.AllWeapons, so none is silently
// excluded from loadouts.
func TestRegistry_CoversEveryWeapon(t *testing.T) {
	weaponsDir := filepath.Join(RepoRoot(t), "internal", "weapon", "weapons")
	entries, err := os.ReadDir(weaponsDir)
	if err != nil {
		t.Fatalf("read weapons dir: %v", err)
	}

	implemented := map[string]bool{}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(weaponsDir, e.Name()), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Recv.NumFields() == 0 || fn.Name.Name != "ID" {
				continue
			}
			if !returnsWeaponID(fn) {
				continue
			}
			if recv, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
				implemented[recv.Name] = true
			}
		}
	}

	registered := map[string]bool{}
	for _, w := range registry.AllWeapons {
		registered[reflect.TypeOf(w).Name()] = true
	}

	for name := range implemented {
		if !registered[name] {
			t.Errorf("weapon %s is implemented but missing from registry.AllWeapons", name)
		}
	}
	for name := range registered {
		if !implemented[name] {
			t.Errorf("registry.AllWeapons lists %s, which is not an implemented weapon type", name)
		}
	}
}

// returnsWeaponID reports whether fn's single result is ids.WeaponID — the signature that
// marks a weapon type apart from its activated-ability card, which returns ids.CardID.
func returnsWeaponID(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}
	sel, ok := fn.Type.Results.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "ids" && sel.Sel.Name == "WeaponID"
}
