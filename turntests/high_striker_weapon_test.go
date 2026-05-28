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

// Tests that High Striker's "next attack hits" rider fires on a follow-up weapon swing.
func TestHighStriker_WeaponHitCreatesCopper(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{testutils.ClubWeapon{}}, nil)
	hand := []card.Card{
		cards.HighStrikerRed{},
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	if got := summary.State.CopperCount(); got != 6 {
		t.Fatalf("Copper at start of next turn = %d, want 6 (HSR rider on Club swing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
