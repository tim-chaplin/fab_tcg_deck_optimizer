package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Tests that AppendHand keeps the hand ordered by Card.ID() regardless of insertion order
// — the canonical-multiset invariant the eval cache depends on. DrawOne shares the same
// insertHandSorted path.
func TestAppendHand_KeepsHandSortedByID(t *testing.T) {
	ge := New()
	for _, id := range []ids.CardID{5, 2, 9, 2, 1, 7} {
		ge.AppendHand(stubCard{id: id})
	}
	got := ge.Hand()
	if len(got) != 6 {
		t.Fatalf("hand len = %d, want 6", len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].ID() > got[i].ID() {
			t.Errorf("hand not sorted at index %d: ID %d before %d", i, got[i-1].ID(), got[i].ID())
		}
	}
}
