package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider drives the chain runner
// through (*Deck).EvalOneTurnForTesting: Aether Slash Red costs 1, Malefic Incantation Blue
// pitches 3 (non-attack action), so the only feasible chain pitches Malefic to play Aether
// Slash. Pitch attribution gives Aether Slash the Malefic pitch, firing the +1 arcane
// rider — Aether Slash deals 4 (its base power) plus 1 (the rider). The Malefic pitch's
// residual carry (3 - 1 = 2) is fine; it just goes unused.
func TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider(t *testing.T) {
	hand := []sim.Card{cards.AetherSlashRed{}, cards.MaleficIncantationBlue{}}
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	state := d.EvalOneTurnForTesting(0, nil, hand)
	if state.PrevTurnValue != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Aether Slash 4 + rider 1)", state.PrevTurnValue)
	}
}

// TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider mirrors the previous test but
// swaps Malefic for testutils.YellowAttack (pitch 2, attack-typed). Now no non-attack
// pitch is available, so the rider can't fire under any ordering. The optimizer pitches
// the Yellow Attack to fund Aether Slash; total damage is just Aether Slash's 4.
func TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider(t *testing.T) {
	hand := []sim.Card{cards.AetherSlashRed{}, testutils.YellowAttack{}}
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	state := d.EvalOneTurnForTesting(0, nil, hand)
	if state.PrevTurnValue != 4 {
		t.Fatalf("PrevTurnValue = %d, want 4 (Aether Slash base power, no rider)", state.PrevTurnValue)
	}
}

// fillerDeck is a no-op deck body for EvalOneTurnForTesting calls that supply their own
// initialHand. The cards never enter play this turn (the caller's hand is the only thing
// the chain runner sees) but EvalOneTurnForTesting still wants a non-empty Deck.Cards so
// the post-turn deal can pull a turn-2 hand without short-circuiting to a zero state.
func fillerDeck() []sim.Card {
	return []sim.Card{
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
	}
}
