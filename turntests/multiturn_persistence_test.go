package turntests

// Multi-turn tests that pin behaviour requiring persistent state to carry across the
// turn boundary correctly. Catches regressions where ResetEphemeralState wipes a
// persistent field, or per-perm cloning fails to inherit one.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// TestEvalTwoTurns_EquippedWeaponPersistsAcrossTurns pins that an equipped weapon stays
// swingable on turn 2 — catches a regression where weapons aren't carried in the cross-turn
// state. Hand1 = 1 BluePitch (pitches 3, funds the 2-cost swing); deck = 1 BluePitch
// (drawn into turn 2 for the same line). Each turn: pitch the Blue, swing Nebula Blade
// (1 base attack + 1 from the on-hit Runechant credit) = 2 damage. Expected:
// turn1.Value == 2 AND turn2.Value == 2. If the weapon goes missing on turn 2 (persistence
// broken), turn2.Value drops to 0.
func TestEvalTwoTurns_EquippedWeaponPersistsAcrossTurns(t *testing.T) {
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{
		testutils.BluePitch{},
	})
	hand := []card.Card{testutils.BluePitch{}}

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, nil, hand)

	if turn1.Value != 2 {
		t.Errorf("turn 1 Value = %d, want 2 (pitch Blue, swing Nebula Blade for 1+1)\nBestLine: %s",
			turn1.Value, formatBestLine(turn1.BestLine))
	}
	if turn2.Value != 2 {
		t.Errorf("turn 2 Value = %d, want 2 (weapon should still be equipped; pitch Blue, swing again for 1+1)\nBestLine: %s",
			turn2.Value, formatBestLine(turn2.BestLine))
	}
}
