package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestWaterTheSeeds_NoAttackReturnsBase: with nothing attack-typed in CardsRemaining the +1 rider
// fizzles and each variant returns its base power.
func TestWaterTheSeeds_NoAttackReturnsBase(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{WaterTheSeedsRed{}, 3},
		{WaterTheSeedsYellow{}, 2},
		{WaterTheSeedsBlue{}, 1},
	}
	for _, tc := range cases {
		s := &card.TurnState{}
		tc.c.Play(s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (no lookahead target)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestWaterTheSeeds_HighPowerFizzles: a power-2 attack is past the base-{p}-<=1 gate, so the
// rider keeps searching. With no matching attack below it, the bonus fizzles.
func TestWaterTheSeeds_HighPowerFizzles(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 2)}}}
	(WaterTheSeedsRed{}).Play(s, &card.CardState{Card: WaterTheSeedsRed{}})
	if got := s.Value; got != 3 {
		t.Errorf("Play() = %d, want 3 (power 2 > 1 → no bonus)", got)
	}
}

// TestWaterTheSeeds_LowPowerTriggersBonus: a power-1 attack matches the gate and fires the
// +1 rider — the buff lands on the target's BonusAttack so its EffectiveAttack picks up
// the +1, not the granter's chain step.
func TestWaterTheSeeds_LowPowerTriggersBonus(t *testing.T) {
	for _, c := range []card.Card{WaterTheSeedsRed{}, WaterTheSeedsYellow{}, WaterTheSeedsBlue{}} {
		target := &card.CardState{Card: stubGenericAttack(0, 1)}
		s := &card.TurnState{CardsRemaining: []*card.CardState{target}}
		c.Play(s, &card.CardState{Card: c})
		if got := target.BonusAttack; got != 1 {
			t.Errorf("%s: target.BonusAttack = %d, want 1 (power-1 target triggers +1)", c.Name(), got)
		}
	}
}

// TestWaterTheSeeds_SkipsPastNonMatchingAttacks: the "next attack with base {p} <=1" trigger
// lasts until a matching attack resolves, so a power-3 attack scheduled before a power-0
// attack shouldn't consume the rider — the +1 lands on the power-0 target.
func TestWaterTheSeeds_SkipsPastNonMatchingAttacks(t *testing.T) {
	skipped := &card.CardState{Card: stubGenericAttack(0, 3)}
	target := &card.CardState{Card: stubGenericAttack(0, 0)}
	s := &card.TurnState{CardsRemaining: []*card.CardState{skipped, target}}
	(WaterTheSeedsRed{}).Play(s, &card.CardState{Card: WaterTheSeedsRed{}})
	if got := skipped.BonusAttack; got != 0 {
		t.Errorf("skipped.BonusAttack = %d, want 0 (power-3 target shouldn't consume rider)", got)
	}
	if got := target.BonusAttack; got != 1 {
		t.Errorf("target.BonusAttack = %d, want 1 (rider lands on the power-0 attack)", got)
	}
}

// TestWaterTheSeeds_NonAttackInRemainingIgnored: only attack-action cards in CardsRemaining count
// as potential triggers.
func TestWaterTheSeeds_NonAttackInRemainingIgnored(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	(WaterTheSeedsRed{}).Play(s, &card.CardState{Card: WaterTheSeedsRed{}})
	if got := s.Value; got != 3 {
		t.Errorf("Play() = %d, want 3 (non-attack ignored)", got)
	}
}
