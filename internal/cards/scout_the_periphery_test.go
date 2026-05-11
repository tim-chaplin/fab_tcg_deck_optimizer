package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestScoutThePeriphery_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestScoutThePeriphery_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []card.Card{ScoutThePeripheryRed{}, ScoutThePeripheryYellow{}, ScoutThePeripheryBlue{}} {
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestScoutThePeriphery_NonAttackInRemainingFizzles: non-attack action (even from arsenal)
// fails the predicate — only attack actions count as the rider's target.
func TestScoutThePeriphery_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAction(), FromArsenal: true}}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestScoutThePeriphery_HandPlayedAttackFizzles: queued attack action that wasn't played from
// arsenal fails the rider's "next attack action card you play from arsenal" target gate.
func TestScoutThePeriphery_HandPlayedAttackFizzles(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAttack(0, 0)}}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (target attack not from arsenal)", got)
	}
}

// Tests that the per-variant bonus lands on the next arsenal-played attack's BonusAttack
// (granter credits 0; bonus rides on the target).
func TestScoutThePeriphery_NextArsenalAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ScoutThePeripheryRed{}, 3},
		{ScoutThePeripheryYellow{}, 2},
		{ScoutThePeripheryBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0), FromArsenal: true}
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: granter credits %d, want 0 (bonus rides on target)", tc.c.Name(), got)
		}
		if got := target.BonusAttack; got != tc.want {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
