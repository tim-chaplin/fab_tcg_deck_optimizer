package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestWaterTheSeeds_NoAttackReturnsBase: with nothing attack-typed in CardsRemaining the +1 rider
// fizzles and each variant returns its base power.
func TestWaterTheSeeds_NoAttackReturnsBase(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.WaterTheSeedsRed{}, 3},
		{cards.WaterTheSeedsYellow{}, 2},
		{cards.WaterTheSeedsBlue{}, 1},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (no lookahead target)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestWaterTheSeeds_HighPowerFizzles: a power-2 attack is past the base-{p}-<=1 gate, so the
// rider keeps searching. With no matching attack below it, the bonus fizzles.
func TestWaterTheSeeds_HighPowerFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAttack(0, 2)}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WaterTheSeedsRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (power 2 > 1 → no bonus)", got)
	}
}

// TestWaterTheSeeds_LowPowerTriggersBonus: a power-1 attack matches the gate and fires the
// +1 rider — the buff lands on the target's BonusAttack so its EffectiveAttack picks up
// the +1, not the granter's chain step.
func TestWaterTheSeeds_LowPowerTriggersBonus(t *testing.T) {
	for _, c := range []card.Card{cards.WaterTheSeedsRed{}, cards.WaterTheSeedsYellow{}, cards.WaterTheSeedsBlue{}} {
		target := &card.CardState{Card: testutils.GenericAttack(0, 1)}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := target.BonusAttack; got != 1 {
			t.Errorf("%s: target.BonusAttack = %d, want 1 (power-1 target triggers +1)", c.Name(), got)
		}
	}
}

// TestWaterTheSeeds_SkipsPastNonMatchingAttacks: the "next attack with base {p} <=1" trigger
// lasts until a matching attack resolves, so a power-3 attack scheduled before a power-0
// attack shouldn't consume the rider — the +1 lands on the power-0 target.
func TestWaterTheSeeds_SkipsPastNonMatchingAttacks(t *testing.T) {
	skipped := &card.CardState{Card: testutils.GenericAttack(0, 3)}
	target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{skipped, target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WaterTheSeedsRed{}})
	if got := skipped.BonusAttack; got != 0 {
		t.Errorf("skipped.BonusAttack = %d, want 0 (power-3 target shouldn't consume rider)", got)
	}
	if got := target.BonusAttack; got != 1 {
		t.Errorf("target.BonusAttack = %d, want 1 (rider lands on the power-0 attack)", got)
	}
}

// TestWaterTheSeeds_NonAttackInRemainingIgnored: a generic action card in CardsRemaining
// doesn't qualify as "your next attack" — the rider walks past it without firing.
func TestWaterTheSeeds_NonAttackInRemainingIgnored(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WaterTheSeedsRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (non-attack ignored)", got)
	}
}

// Tests that the +1 rider lands on a weapon swing target ("your next attack" has no action-card
// qualifier).
func TestWaterTheSeeds_BonusLandsOnWeaponSwing(t *testing.T) {
	target := &card.CardState{Card: testutils.RunebladeWeapon{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WaterTheSeedsRed{}})
	if got := target.BonusAttack; got != 1 {
		t.Errorf("weapon BonusAttack = %d, want 1 (no 'action card' qualifier on 'next attack')", got)
	}
}
