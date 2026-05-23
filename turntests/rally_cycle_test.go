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
// beats keeping the card — discarding a weak spare, keeping a strong one to attack, and
// holding an unplayable spare (into the arsenal) when the base block already suffices.
func TestRallyCycle_DiscardsHeldCardOnlyWhenWorthIt(t *testing.T) {
	cases := []struct {
		name     string
		incoming int
		spare    card.Card
		want     int
		wantHeld bool
	}{
		// Discarding the 2-power spare for +3{d} blocks all 5 (value 5); playing it instead
		// blocks 2 and deals 2 (value 4). The discard wins.
		{"discards a weak spare", 5, testutils.GenericAttack(0, 2), 5, false},
		// Keeping the 4-power spare blocks 2 and deals 4 (value 6); discarding it blocks 5
		// (value 5). Keeping it to attack wins.
		{"keeps a strong spare to attack", 5, testutils.GenericAttack(0, 4), 6, false},
		// A cost-3 spare can't be funded, and the base 2 block already covers the 2 incoming.
		// The +3{d} would be dead weight, so the spare is held — promoted to the arsenal.
		{"holds an unplayable spare", 2, testutils.GenericAttack(3, 2), 2, true},
	}
	for _, rally := range []card.Card{cards.RallyTheCoastGuardRed{}, cards.RallyTheRearguardRed{}} {
		for _, c := range cases {
			d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
			prior := gameengine.GameStateBuilder().SetIncomingDamage(c.incoming).Build()
			summary := sim.EvalOneTurnForTesting(d, prior, []card.Card{rally, c.spare})

			if got := summary.Value; got != c.want {
				t.Errorf("%s / %s: Value = %d, want %d\nBestLine: %s",
					rally.Name(), c.name, got, c.want, formatBestLine(summary.BestLine))
			}
			// A held spare is promoted to the arsenal; a discarded or played one is not.
			if c.wantHeld && summary.State.Arsenal() != c.spare {
				t.Errorf("%s / %s: spare not held — arsenal = %v, want the spare\nBestLine: %s",
					rally.Name(), c.name, summary.State.Arsenal(), formatBestLine(summary.BestLine))
			}
		}
	}
}
