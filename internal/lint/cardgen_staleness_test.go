package lint

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cardgen"
)

// Tests that every cardgen-generated file is current: the committed _gen.go files in the
// cards package and its notimplemented / unplayable subpackages, plus the registry's
// cards_gen.go, match what cardgen produces from the yamls today. Fails when a card yaml
// was added or changed without rerunning `go generate ./...`.
func TestCardgen_GeneratedFilesAreCurrent(t *testing.T) {
	root := RepoRoot(t)
	cardsDir := filepath.Join(root, "internal", "card", "cards")
	dirs := []string{
		cardsDir,
		filepath.Join(cardsDir, "notimplemented"),
		filepath.Join(cardsDir, "unplayable"),
	}
	registryPath := filepath.Join(root, "internal", "registry", "cards_gen.go")

	files, err := cardgen.Generate(dirs, registryPath)
	if err != nil {
		t.Fatalf("cardgen.Generate: %v", err)
	}

	for path, want := range files {
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		got, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: %v — run `go generate ./...`", rel, err)
			continue
		}
		// Committed blobs are LF; the working tree may be CRLF. Compare EOL-agnostic.
		if normalizeEOL(got) != normalizeEOL(want) {
			t.Errorf("%s is stale — run `go generate ./...`", rel)
		}
	}
}

// normalizeEOL collapses CRLF to LF so a CRLF working-tree checkout compares equal to the
// LF-formatted bytes cardgen produces.
func normalizeEOL(b []byte) string {
	return string(bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")))
}
