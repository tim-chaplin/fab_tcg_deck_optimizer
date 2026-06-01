package registry

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
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

func TestRanking_PromoteMovesToTop(t *testing.T) {
	r := NewRanking()
	promoted := r.order[5]
	above := append([]CardID(nil), r.order[:5]...) // cards currently ranked above it
	r.PromoteAll([]CardID{promoted})
	if got := r.Rank(promoted); got != 0 {
		t.Errorf("promoted card rank = %d, want 0", got)
	}
	for i, id := range above { // everything above it slides down exactly one slot
		if got := r.Rank(id); got != i+1 {
			t.Errorf("card formerly at %d now at %d, want %d", i, got, i+1)
		}
	}
}

// TestRanking_PromoteAllPartitions: a set of scattered cards leads the order in their current
// relative order, the rest follow in theirs (a stable partition).
func TestRanking_PromoteAllPartitions(t *testing.T) {
	r := NewRanking()
	a, b, c := r.order[7], r.order[2], r.order[9] // old ranks 7, 2, 9
	formerTop := r.order[0]                       // not promoted
	r.PromoteAll([]CardID{a, b, c})
	// Promoted cards lead, kept in their old relative order: b(2) < a(7) < c(9).
	if r.Rank(b) != 0 || r.Rank(a) != 1 || r.Rank(c) != 2 {
		t.Errorf("promoted set not in relative order: b=%d a=%d c=%d, want 0,1,2", r.Rank(b), r.Rank(a), r.Rank(c))
	}
	// The card formerly at rank 0 is now the first of the rest — just behind the 3 promoted.
	if got := r.Rank(formerTop); got != 3 {
		t.Errorf("former top now at %d, want 3 (just behind the promoted set)", got)
	}
}

func TestRanking_PromoteTopAndUnknownAreNoops(t *testing.T) {
	r := NewRanking()
	before := append([]CardID(nil), r.order...)
	absent := CardID(0)
	for _, id := range r.order {
		if id > absent {
			absent = id
		}
	}
	absent++                           // one past the largest pool ID — not in the ranking
	r.PromoteAll([]CardID{r.order[0]}) // already best
	r.PromoteAll([]CardID{absent})     // not in the ranking
	r.PromoteAll(nil)                  // empty list
	for i, id := range before {
		if r.Rank(id) != i {
			t.Errorf("no-op PromoteAll moved card %d to %d", i, r.Rank(id))
		}
	}
}

// TestRanking_PromoteConcurrent checks concurrent PromoteAll / Rank keeps order a valid permutation
// with pos consistent (run with -race).
func TestRanking_PromoteConcurrent(t *testing.T) {
	r := NewRanking()
	n := len(r.order)
	ids := append([]CardID(nil), r.order...) // snapshot the pool IDs

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(seed))
			for i := 0; i < 2000; i++ {
				r.PromoteAll([]CardID{ids[rng.Intn(n)], ids[rng.Intn(n)]})
				_ = r.Rank(ids[rng.Intn(n)])
			}
		}(int64(g + 1))
	}
	wg.Wait()

	if len(r.order) != n {
		t.Fatalf("order length changed to %d, want %d", len(r.order), n)
	}
	seen := make(map[CardID]bool, n)
	for i, id := range r.order {
		if seen[id] {
			t.Fatalf("duplicate card at index %d", i)
		}
		seen[id] = true
		if r.pos[id] != i {
			t.Errorf("pos/order inconsistent: pos says %d, order index is %d", r.pos[id], i)
		}
	}
	if len(seen) != n {
		t.Errorf("permutation lost cards: %d distinct, want %d", len(seen), n)
	}
}

func TestRanking_SaveLoadRoundTrip(t *testing.T) {
	r := NewRanking()
	r.PromoteAll([]CardID{r.order[10]}) // perturb the order
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
