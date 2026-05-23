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

// Tests that Bluster Buff picks mode 1 (pay extra {r} for full 6{p}) when the partition's
// pitch supply has the resource to spare.
func TestModalCost_BlusterBuffPicksMode1WhenAffordable(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.BlusterBuffRed{},
		testutils.FakeBlueResource(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 6 {
		t.Fatalf("Value = %d, want 6 (mode 1: pitch BluePitch for 3{r}, pay 2{r}, attack 6{p})", got)
	}
}

// Tests that Bluster Buff falls back to mode 0 (5{p}, no extra cost) when only the
// printed 1{r} is available.
func TestModalCost_BlusterBuffPicksMode0WhenBudgetTight(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.BlusterBuffRed{},
		cards.BlusterBuffRed{},
	}
	// Pitch one Bluster Buff (1{r}), play the other for cost 1{r} mode 0 (5{p}). Mode 1
	// would need 2{r} which the 1-pitch supply can't fund.
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (mode 0: pitch BB for 1{r}, attack 5{p})", got)
	}
}

// Tests that Look Tuff picks mode 1 when the partition has 4{r} of pitch (printed 3 + extra 1).
func TestModalCost_LookTuffPicksMode1WhenAffordable(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.LookTuffRed{},
		testutils.FakeBlueResource(),
		cards.LookTuffRed{}, // pitch supply: BluePitch 3 + 1 = 4
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (mode 1: pay 4{r}, attack 8{p})", got)
	}
}
