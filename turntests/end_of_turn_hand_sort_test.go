package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the post-turn hand is sorted by Card.ID() even when end-of-turn refill draws
// zero cards (held count already matches intellect, OR deck is empty). The attack-turn runner's
// RemoveFromHand uses swap-with-last so the held cards survive in non-canonical order;
// the cache key + attack-turn runner both rely on the hand being sorted at findBest entry, so
// playOneTurn must always sort the post-turn hand. The pre-fix code only sorted when
// toDraw > 0, leaving the unsorted held order to leak into next turn's findBest call.
func TestEndOfTurnHand_SortedEvenWithZeroDraws(t *testing.T) {
	// Trade In Red (ID 486) is the only attack action — it plays, and its Play discards
	// a Held card (via ge.Discard, which scans for the first Held slot, index 0). Discard
	// is a swap-with-last pop, so removing the lowest-ID held card brings the
	// highest-ID held card to index 0, leaving held = [OasisRespite (319), BrushOff (153)]
	// — UNSORTED. Intel=2 with 2 held survivors makes toDraw == 0, the empty deck removes
	// any mid-attack-turn draw confound, and the pre-fix code skipped the sort on that path.
	// Hand of 5: TradeInRed is the attack, two arsenal-eligible DRs (one becomes the
	// post-hoc arsenal pick), and two Resource-typed cards that are arsenal-ineligible so
	// they stay held. After Trade In's Discard pops the first held (ReduceToRunechant at
	// index 0), swap-with-last brings TitaniumBauble to the front: held = [482, 175] —
	// UNSORTED. Intel=2 with 2 held survivors triggers toDraw == 0.
	hand := []card.Card{
		cards.TradeInRed{},            // ID 486 — the attack action
		cards.ReduceToRunechantBlue{}, // ID 60  — DR, discarded by Trade In
		cards.BrushOffBlue{},          // ID 153 — DR, gets post-hoc promoted to arsenal
		cards.CrackedBaubleYellow{},   // ID 175 — Resource, arsenal-ineligible, held
		cards.TitaniumBaubleBlue{},    // ID 482 — Resource, arsenal-ineligible, held
	}
	d := deck.New(testutils.Hero{Intel: 2}, nil, nil)

	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	got := summary.State.Hand()
	for i := 1; i < len(got); i++ {
		if got[i-1].ID() > got[i].ID() {
			ids := make([]int, len(got))
			for j, c := range got {
				ids[j] = int(c.ID())
			}
			t.Fatalf("post-turn hand not sorted by ID at index %d: %v (held cards must be canonicalised before becoming next turn's input to findBest)", i, ids)
		}
	}
}
