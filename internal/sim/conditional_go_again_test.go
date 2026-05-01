package sim_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestConditionalGoAgainMarkerCoverage uses static analysis to enforce the
// ConditionalGoAgain marker contract. For every non-test .go file in internal/cards,
// the test parses the file's AST, looks for any assignment whose LHS is
// self.GrantedGoAgain (the convention by which a card grants Go again to itself in
// Play), and — when a flip is found — asserts that every card variant struct declared
// in that file implements the sim.ConditionalGoAgain marker.
//
// Static analysis avoids the upkeep cost of a permissive runtime probe (no growing
// list of "trigger flags this archetype gates on"); the parser sees every code path,
// regardless of the runtime conditions that would actually fire it.
//
// Conventions this lint relies on:
//   - The parameter that holds the playing card's own CardState is named "self".
//     Grants to other cards in CardsRemaining use a different name (commonly "pc")
//     and are correctly skipped — those don't make the granting card itself
//     conditionally go-again.
//   - The helper that flips self.GrantedGoAgain lives in the same file as the
//     variant struct declarations. A helper in a different file would slip past this
//     lint; that's the trade-off vs full call-graph analysis. New cards should
//     follow the in-file convention.
func TestConditionalGoAgainMarkerCoverage(t *testing.T) {
	cardByTypeName := buildCardByTypeName()

	// Tests run with CWD == the package source directory. internal/sim → internal/cards.
	files, err := filepath.Glob("../cards/*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Errorf("parse %s: %v", path, err)
			continue
		}
		if !fileFlipsSelfGrantedGoAgain(f) {
			continue
		}
		base := filepath.Base(path)
		for _, name := range topLevelTypeNames(f) {
			c, ok := cardByTypeName[name]
			if !ok {
				continue // not a registered card (helper struct, base type, etc.)
			}
			if _, ok := c.(ConditionalGoAgain); !ok {
				t.Errorf("%s: %s mutates self.GrantedGoAgain in this file but does not implement sim.ConditionalGoAgain — add `ConditionalGoAgain() {}`",
					base, name)
			}
		}
	}
}

// buildCardByTypeName indexes every registered card by its concrete Go type name (e.g.
// "RuneragerSwarmRed") so the static-analysis loop can match a struct name in the AST
// to its Card instance.
func buildCardByTypeName() map[string]Card {
	m := map[string]Card{}
	for _, id := range DeckableCards() {
		c := GetCard(id)
		if c == nil {
			continue
		}
		m[reflect.TypeOf(c).Name()] = c
	}
	return m
}

// fileFlipsSelfGrantedGoAgain reports whether f contains an assignment whose LHS is
// self.GrantedGoAgain. Grants to OTHER cards via CardsRemaining iteration (e.g.,
// Mauvrion Skies's pc.GrantedGoAgain) are deliberately excluded; those don't make
// the granting card itself conditionally go-again.
func fileFlipsSelfGrantedGoAgain(f *ast.File) bool {
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if found {
			return false
		}
		a, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, lhs := range a.Lhs {
			sel, ok := lhs.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if sel.Sel.Name != "GrantedGoAgain" {
				continue
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == "self" {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// topLevelTypeNames returns the names of every type declared at file scope. Used by
// the static analysis to find card variant structs that need the marker.
func topLevelTypeNames(f *ast.File) []string {
	var names []string
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			names = append(names, ts.Name.Name)
		}
	}
	return names
}
