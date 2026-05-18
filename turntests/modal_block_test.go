package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Brothers in Arms picks mode 1 (pay 1{r} for +2{d}) when the partition's
// defense pitch supply has the spare resource — Toughen Up's printed 2{r} cost plus BIA's
// extra {r} fits the 3{r} from the Blue Pitch.
func TestModalBlock_BrothersInArmsPicksMode1WhenAffordable(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{
		cards.BrothersInArmsRed{},
		cards.ToughenUpBlue{},
		testutils.BluePitch{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(10).Build(), hand)
	got := summary.Value
	// Toughen Up DR: 4{d}; BIA mode 1: 2 + 2 = 4{d}; pitch supply 3{r} covers 2 (TU) + 1 (BIA mode 1).
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (Toughen Up 4 + BIA mode 1: 2+2 bonus)", got)
	}
}

// Tests that Brothers in Arms falls back to mode 0 when no spare {r} is available — the
// hand has no pitch source besides BIA itself, and pitching BIA would forfeit the block.
func TestModalBlock_BrothersInArmsFallsBackToMode0(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{cards.BrothersInArmsRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(10).Build(), hand)
	got := summary.Value
	if got != 2 {
		t.Fatalf("Value = %d, want 2 (BIA mode 0: printed 2{d}, no spare {r} for mode 1)", got)
	}
}
