package lint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Tests that every card file calling GameEngine.AddResourcePoints also declares
// MaxResourcePoints, the card.ResourceSource contract the attack-budget prune relies on.
func TestResourceSource_AddResourcePointsImpliesMaxResourcePoints(t *testing.T) {
	cardsRoot := filepath.Join(RepoRoot(t), "internal", "card", "cards")
	err := filepath.WalkDir(cardsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_gen.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		callsAdd, declaresMax := false, false
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if x.Sel.Name == "AddResourcePoints" {
					callsAdd = true
				}
			case *ast.FuncDecl:
				if x.Recv != nil && x.Name.Name == "MaxResourcePoints" {
					declaresMax = true
				}
			}
			return true
		})
		if callsAdd && !declaresMax {
			t.Errorf("%s: reaches GameEngine.AddResourcePoints but declares no "+
				"MaxResourcePoints — implement card.ResourceSource so the attack-budget "+
				"prune accounts for the resource points it adds.", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}
