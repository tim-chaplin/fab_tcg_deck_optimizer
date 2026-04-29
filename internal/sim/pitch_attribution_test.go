package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider drives the chain runner
// end-to-end: Aether Slash Red costs 1, Malefic Incantation Blue pitches 3 (non-attack
// action), so the only feasible chain pitches Malefic to play Aether Slash. Pitch
// attribution gives Aether Slash the Malefic pitch, firing the +1 arcane rider — Aether
// Slash deals 4 (its base power) plus 1 (the rider). The Malefic pitch's residual carry
// (3 - 1 = 2) is fine; it just goes unused.
func TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider(t *testing.T) {
	h := []Card{cards.AetherSlashRed{}, cards.MaleficIncantationBlue{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5 (Aether Slash 4 + rider 1), got %d", got.Value)
	}
}

// TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider mirrors the previous test but
// swaps Malefic for testutils.YellowAttack (pitch 2, attack-typed). Now no non-attack
// pitch is available, so the rider can't fire under any ordering. Best plays the chain
// that pitches the Yellow Attack to fund Aether Slash; total damage is just Aether Slash's
// 4 (no rider).
func TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider(t *testing.T) {
	h := []Card{cards.AetherSlashRed{}, testutils.YellowAttack{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4 (Aether Slash base power, no rider), got %d", got.Value)
	}
}
