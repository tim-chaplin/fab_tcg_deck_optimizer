package gameengine

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/lint"
)

// TestGameStateBuilder_MultiLineWhenMultiSetter enforces the codebase convention that a
// fluent GameStateBuilder() chain carrying two or more Set* calls is broken one call per
// line:
//
//	gameengine.GameStateBuilder().
//		SetHero(h).
//		SetArsenal(a).
//		Build()
//
// Single-setter chains (GameStateBuilder().SetHero(h).Build()) may stay inline. The test
// walks every .go file in the repo so the rule holds across packages.
func TestGameStateBuilder_MultiLineWhenMultiSetter(t *testing.T) {
	root := lint.RepoRoot(t)
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != root && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			lines, ok := gameStateBuilderChainLines(fset, call)
			if !ok {
				return true
			}
			// lines holds [Build, SetN, ..., Set1, GameStateBuilder]; the setter count is
			// everything except the terminal Build and the root GameStateBuilder.
			if len(lines)-2 < 2 {
				return true
			}
			seen := make(map[int]bool, len(lines))
			for _, ln := range lines {
				if seen[ln] {
					pos := fset.Position(call.Lparen)
					t.Errorf("%s: GameStateBuilder() chain with %d Set* calls must be broken one call per line",
						pos, len(lines)-2)
					break
				}
				seen[ln] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}

// gameStateBuilderChainLines reports whether call is the terminal .Build() of a
// GameStateBuilder() chain. When it is, it returns the source line of every method
// selector in the chain ordered outer→inner ([Build, SetN, ..., Set1, GameStateBuilder]);
// otherwise ok is false. Nested Set* / GameStateBuilder calls return ok=false so a chain
// is reported once, from its outermost Build().
func gameStateBuilderChainLines(fset *token.FileSet, call *ast.CallExpr) (lines []int, ok bool) {
	sel, isSel := call.Fun.(*ast.SelectorExpr)
	if !isSel || sel.Sel.Name != "Build" {
		return nil, false
	}
	cur := ast.Expr(call)
	for {
		c, isCall := cur.(*ast.CallExpr)
		if !isCall {
			return nil, false
		}
		switch fn := c.Fun.(type) {
		case *ast.SelectorExpr:
			lines = append(lines, fset.Position(fn.Sel.NamePos).Line)
			if fn.Sel.Name == "GameStateBuilder" {
				return lines, true
			}
			cur = fn.X
		case *ast.Ident:
			// Unqualified GameStateBuilder() — calls from inside package gameengine.
			if fn.Name == "GameStateBuilder" {
				lines = append(lines, fset.Position(fn.NamePos).Line)
				return lines, true
			}
			return nil, false
		default:
			return nil, false
		}
	}
}
