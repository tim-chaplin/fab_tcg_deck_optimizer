package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that two equipped 1H copies of the same weapon each bump their own counter object
// when both swing in one turn (the swing has go again, so one action point chains both).
func TestDuplicateWeapon_CounterAttributedToSwungObject(t *testing.T) {
	d := deck.New(heroes.Viserai,
		[]deck.Weapon{testutils.CounterWeapon{}, testutils.CounterWeapon{}}, nil)
	hand := []card.Card{
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
	}

	summary := sim.EvalOneTurnForTesting(d,
		gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	weapons := summary.State.Weapons()
	if len(weapons) != 2 {
		t.Fatalf("equipped weapons = %d, want 2\nBestLine: %s",
			len(weapons), formatBestLine(summary.BestLine))
	}
	c0, c1 := weapons[0].Count(), weapons[1].Count()
	// Both swing (go-again chains them), so each object holds exactly one counter.
	if c0 != 1 || c1 != 1 {
		t.Errorf("counters = (c0=%d, c1=%d), want (1, 1): each swung object should bump its own counter\nBestLine: %s",
			c0, c1, formatBestLine(summary.BestLine))
	}
}
