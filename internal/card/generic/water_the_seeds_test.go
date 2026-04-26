package generic

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
	if got := s.Value; got != 3{
		t.Errorf("Play() = %d, want 3 (power 2 > 1 → no bonus)", got)
	}
}

// TestWaterTheSeeds_LowPowerTriggersBonus: a power-1 attack matches the gate and fires the
// +1 rider, so each variant returns base + 1.
func TestWaterTheSeeds_LowPowerTriggersBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{WaterTheSeedsRed{}, 4},
		{WaterTheSeedsYellow{}, 3},
		{WaterTheSeedsBlue{}, 2},
	}
	for _, tc := range cases {
		s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 1)}}}
		tc.c.Play(s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (power-1 target triggers +1)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestWaterTheSeeds_SkipsPastNonMatchingAttacks: the "next attack with base {p} <=1" trigger lasts
// until a matching attack resolves, so a power-3 attack scheduled before a power-0 attack shouldn't
// consume the rider.
func TestWaterTheSeeds_SkipsPastNonMatchingAttacks(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.CardState{
		{Card: stubGenericAttack(0, 3)},
		{Card: stubGenericAttack(0, 0)},
	}}
	(WaterTheSeedsRed{}).Play(s, &card.CardState{Card: WaterTheSeedsRed{}})
	if got := s.Value; got != 4{
		t.Errorf("Play() = %d, want 4 (rider waits for the power-0 attack)", got)
	}
}

// TestWaterTheSeeds_NonAttackInRemainingIgnored: only attack-action cards in CardsRemaining count
// as potential triggers.
func TestWaterTheSeeds_NonAttackInRemainingIgnored(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	(WaterTheSeedsRed{}).Play(s, &card.CardState{Card: WaterTheSeedsRed{}})
	if got := s.Value; got != 3{
		t.Errorf("Play() = %d, want 3 (non-attack ignored)", got)
	}
}
