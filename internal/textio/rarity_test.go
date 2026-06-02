package textio

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadRarityLegal checks the join: a card is Silver-Age-rarity-legal iff it has at least one
// Basic/Common/Rare printing, keyed by Card Unique ID.
func TestLoadRarityLegal(t *testing.T) {
	dir := t.TempDir()
	rarity := filepath.Join(dir, "rarity.csv")
	printing := filepath.Join(dir, "card-printing.csv")
	write := func(path, content string) {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(rarity, "Shorthand\tText\nC\tCommon\nR\tRare\nM\tMajestic\nB\tBasic\n")
	// A: Common (legal). B: only Majestic (illegal). C: Majestic + Rare (legal). D: Basic (legal).
	write(printing, "Card Unique ID\tRarity\nA\tC\nB\tM\nC\tM\nC\tR\nD\tB\n")

	legal, err := LoadRarityLegal(printing, rarity)
	if err != nil {
		t.Fatalf("LoadRarityLegal: %v", err)
	}
	for id, want := range map[string]bool{"A": true, "B": false, "C": true, "D": true} {
		if legal[id] != want {
			t.Errorf("legal[%q] = %v, want %v", id, legal[id], want)
		}
	}
}
