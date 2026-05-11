package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestSmashingGoodTime_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSmashingGoodTime_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSmashingGoodTime_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSmashingGoodTime_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAction()}}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: SmashingGoodTimeRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSmashingGoodTime_NextAttackGrantsBonusAttack: arsenal-played copy with a queued
// attack action grants the per-variant bonus (Red +3, Yellow +2, Blue +1) onto the target's
// BonusAttack. Granter returns 0; the +N attributes to the buffed attack.
func TestSmashingGoodTime_NextAttackGrantsBonusAttack(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{SmashingGoodTimeRed{}, 3},
		{SmashingGoodTimeYellow{}, 2},
		{SmashingGoodTimeBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
		self := &card.CardState{Card: tc.c, FromArsenal: true}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// TestSmashingGoodTime_HandPlayedFizzles: hand-played copy fails the from-arsenal gate.
func TestSmashingGoodTime_HandPlayedFizzles(t *testing.T) {
	for _, c := range []sim.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAttack(0, 0)}}})
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (hand-played)", c.Name(), got)
		}
	}
}
