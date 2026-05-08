package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that High Striker's "next attack hits" rider fires on a follow-up weapon swing.
func TestHighStriker_WeaponHitCreatesCopper(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{testutils.ClubWeapon{}}, fillerDeck())
	hand := []sim.Card{
		cards.HighStrikerRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	summary := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	state := summary.State
	if got := state.Copper(); got != 6 {
		t.Fatalf("Copper at start of next turn = %d, want 6 (HSR rider on Club swing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
