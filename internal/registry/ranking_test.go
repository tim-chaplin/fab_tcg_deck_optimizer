package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRanking_SeedsWholePool(t *testing.T) {
	r := NewRanking()
	legal := (Registry{}).LegalCards()
	if len(r.order) != len(legal) {
		t.Fatalf("ranking has %d cards, pool has %d", len(r.order), len(legal))
	}
	seen := make(map[int]bool, len(legal))
	for _, c := range legal {
		rank := r.Rank(c.ID())
		if rank < 0 || rank >= len(legal) {
			t.Errorf("card %s rank %d out of [0,%d)", c.DisplayName(), rank, len(legal))
		}
		if seen[rank] {
			t.Errorf("duplicate rank %d", rank)
		}
		seen[rank] = true
	}
}

func TestRanking_RecordResultPromotesWinner(t *testing.T) {
	r := NewRanking()
	best, challenger := r.order[0], r.order[5]
	r.RecordResult(challenger, best) // lower-ranked card wins → they swap slots
	if got := r.Rank(challenger); got != 0 {
		t.Errorf("winner rank = %d, want 0 (took the loser's slot)", got)
	}
	if got := r.Rank(best); got != 5 {
		t.Errorf("loser rank = %d, want 5 (took the winner's slot)", got)
	}
}

func TestRanking_RecordResultNoopWhenWinnerAlreadyAbove(t *testing.T) {
	r := NewRanking()
	top, low := r.order[0], r.order[5]
	r.RecordResult(top, low) // winner already ranked above loser
	if r.Rank(top) != 0 || r.Rank(low) != 5 {
		t.Errorf("ranks changed on a no-op result: top=%d low=%d", r.Rank(top), r.Rank(low))
	}
}

func TestRanking_SaveLoadRoundTrip(t *testing.T) {
	r := NewRanking()
	r.RecordResult(r.order[10], r.order[2]) // perturb the order
	path := filepath.Join(t.TempDir(), "ranking.json")
	if err := r.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := LoadRanking(path)
	if err != nil {
		t.Fatalf("LoadRanking: %v", err)
	}
	if len(loaded.order) != len(r.order) {
		t.Fatalf("loaded %d cards, saved %d", len(loaded.order), len(r.order))
	}
	for i := range r.order {
		if loaded.order[i] != r.order[i] {
			t.Fatalf("order differs at index %d", i)
		}
	}
}

func TestLoadRanking_MissingFileSeedsPool(t *testing.T) {
	r, err := LoadRanking(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err != nil {
		t.Fatalf("LoadRanking: %v", err)
	}
	if len(r.order) != len((Registry{}).LegalCards()) {
		t.Errorf("missing-file ranking has %d cards, want full pool", len(r.order))
	}
}

func TestLoadRanking_PartialFileTopsListedAppendsRest(t *testing.T) {
	full := NewRanking()
	nameA := GetCard(full.order[20]).DisplayName()
	nameB := GetCard(full.order[3]).DisplayName()
	// Listed cards lead, an unknown name is dropped, and the rest of the pool fills in below.
	data, _ := json.Marshal([]string{nameA, "Not A Real Card [Z]", nameB})
	path := filepath.Join(t.TempDir(), "partial.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRanking(path)
	if err != nil {
		t.Fatalf("LoadRanking: %v", err)
	}
	if r.Rank(full.order[20]) != 0 || r.Rank(full.order[3]) != 1 {
		t.Errorf("listed cards not at top: %d, %d", r.Rank(full.order[20]), r.Rank(full.order[3]))
	}
	if len(r.order) != len((Registry{}).LegalCards()) {
		t.Errorf("ranking size %d, want full pool (missing cards appended)", len(r.order))
	}
}
