package card

// Static lint check enforcing that no card-effect file references the framework-only
// TurnState mutation backdoors. Go's package model has no "framework-only" visibility —
// SetDeck / SetGraveyard are exported on TurnState because the framework (internal/hand,
// internal/deck) needs them, which means card subpackages can also call them. This test
// closes the loophole at build time: it walks every .go file under internal/card/...
// subdirectories (any depth — covers future subpackage additions), parses each AST, and
// fails the test if any references a forbidden identifier. CI catches the bypass before
// merge.
//
// The walk EXCLUDES the internal/card package itself (where the framework methods are
// legitimately defined and tested) — only its subdirectories are scanned, since those
// are where card-effect files live.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// frameworkOnlyTurnStateMethods is the deny-list. Cards that need to mutate the deck or
// graveyard route through PopDeckTop / PrependToDeck / TutorFromDeck / BanishFromGraveyard
// / AddToGraveyard — all of which flip IsCacheable when reading or write append-only.
// SetDeck / SetGraveyard wholesale-replace the slice without flipping, so a card calling
// them would silently bypass the cacheable-tracking guarantee.
var frameworkOnlyTurnStateMethods = map[string]struct{}{
	"SetDeck":      {},
	"SetGraveyard": {},
}

// TestNoFrameworkOnlySettersInCardSubpackages recursively walks every .go file under
// internal/card/<subdir>/... (all depths) and fails if any references a deny-listed
// identifier. Includes _test.go files — card-test files seeding state should also use
// NewTurnState + safe verbs, not raw SetDeck/SetGraveyard. New card-effect subpackages
// added later are picked up automatically.
func TestNoFrameworkOnlySettersInCardSubpackages(t *testing.T) {
	root := "." // tests run from the package's directory: internal/card
	scanned := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip the package's own directory — only subdirectories host card-effect
			// files. (Files at the root are framework-allowed: turn_state.go etc.)
			if path == root {
				return nil
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Skip files that live directly under the card package root. Only subdirectory
		// files are card-effect files.
		if filepath.Dir(path) == root {
			return nil
		}
		scanFileForForbiddenRefs(t, path)
		scanned++
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	if scanned == 0 {
		t.Fatalf("no card-subpackage .go files scanned — layout changed or test running from wrong cwd")
	}
}

// scanFileForForbiddenRefs parses path and reports a test failure for each AST selector
// expression whose Sel name is in the deny-list. The receiver type doesn't need to be
// resolved — *card.TurnState is the only type in the codebase exposing these names, so a
// bare SelectorExpr `<anything>.SetDeck` is conclusive. The lint err's on the side of
// over-reporting; if a future unrelated symbol ever shadows one of these names, deny-list
// it explicitly here.
func scanFileForForbiddenRefs(t *testing.T, path string) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Errorf("parse %s: %v", path, err)
		return
	}
	ast.Inspect(f, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		name := sel.Sel.Name
		if _, banned := frameworkOnlyTurnStateMethods[name]; !banned {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		recipient := "<chained>"
		if ok {
			recipient = ident.Name
		}
		t.Errorf("%s:%s: card subpackage references framework-only TurnState method %q (via %s) — use the safe verbs (PopDeckTop / PrependToDeck / TutorFromDeck / BanishFromGraveyard / AddToGraveyard) instead",
			path, fset.Position(sel.Pos()), name, recipient)
		return true
	})
}
