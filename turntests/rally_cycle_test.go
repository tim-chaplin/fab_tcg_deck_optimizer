package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a defending Rally discards a Held card for the printed +3{d} only when that
// beats keeping the card — winning over a weak spare, losing to a strong one, and skipped
// when the base block already covers the incoming damage.
func TestRallyCycle_DiscardsHeldCardOnlyWhenWorthIt(t *testing.T) {
	cases := []struct {
		name     string
		incoming int
		spare    card.Card
		want     int
	}{
		// Discarding the 2-power spare for +3{d} blocks all 5 (value 5); playing it instead
		// blocks 2 and deals 2 (value 4). The discard wins.
		{"discards a weak spare", 5, testutils.GenericAttack(0, 2), 5},
		// Keeping the 4-power spare blocks 2 and deals 4 (value 6); discarding it blocks 5
		// (value 5). Keeping it wins.
		{"keeps a strong spare", 5, testutils.GenericAttack(0, 4), 6},
		// The base 2 block already covers the 2 incoming, so the spare is kept to attack:
		// block 2 plus 2 damage (value 4) beats a pointless discard (block 2, value 2).
		{"no discard when the base block suffices", 2, testutils.GenericAttack(0, 2), 4},
	}
	for _, rally := range []card.Card{cards.RallyTheCoastGuardRed{}, cards.RallyTheRearguardRed{}} {
		for _, c := range cases {
			d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
			prior := gameengine.GameStateBuilder().SetIncomingDamage(c.incoming).Build()
			summary := sim.EvalOneTurnForTesting(d, prior, []card.Card{rally, c.spare})
			if got := summary.Value; got != c.want {
				t.Errorf("%s / %s: Value = %d, want %d\nBestLine: %s",
					rally.Name(), c.name, got, c.want, formatBestLine(summary.BestLine))
			}
		}
	}
}
