package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Pins High Striker's "next time an attack you control hits" rider firing on a weapon
// swing — High Striker plays go-again, the 1-power Club swing then lands inside the
// LikelyToHit window and the trigger creates the per-pitch Copper count.
func TestHighStriker_WeaponHitCreatesCopper(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{testutils.ClubWeapon{}}, fillerDeck())
	hand := []sim.Card{
		cards.HighStrikerRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if got := state.Copper(); got != 6 {
		t.Fatalf("Copper at start of next turn = %d, want 6 (HSR rider on Club swing)\nBestLine: %s",
			got, formatBestLine(state.BestLine))
	}
}
