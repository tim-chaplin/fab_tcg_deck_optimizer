package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that High Striker's "next attack hits" rider fires on a follow-up weapon swing.
func TestHighStriker_WeaponHitCreatesCopper(t *testing.T) {
	d := deck.New(heroes.Viserai{}, []deck.Weapon{testutils.ClubWeapon{}}, fillerDeck())
	hand := []deck.Card{
		cards.HighStrikerRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if got := state.CopperCount(); got != 6 {
		t.Fatalf("Copper at start of next turn = %d, want 6 (HSR rider on Club swing)\nBestLine: %s",
			got, formatBestLine(state.BestLine))
	}
}
