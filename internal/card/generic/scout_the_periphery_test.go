package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestScoutThePeriphery_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestScoutThePeriphery_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{ScoutThePeripheryRed{}, ScoutThePeripheryYellow{}, ScoutThePeripheryBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestScoutThePeriphery_NonAttackInRemainingFizzles: non-attack action (even from arsenal)
// fails the predicate — only attack actions count as the rider's target.
func TestScoutThePeriphery_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction(), FromArsenal: true}}}
	(ScoutThePeripheryRed{}).Play(&s, &card.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value; got != 0{
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestScoutThePeriphery_HandPlayedAttackFizzles: queued attack action that wasn't played from
// arsenal fails the rider's "next attack action card you play from arsenal" target gate.
func TestScoutThePeriphery_HandPlayedAttackFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 0)}}}
	(ScoutThePeripheryRed{}).Play(&s, &card.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value; got != 0{
		t.Errorf("Play() = %d, want 0 (target attack not from arsenal)", got)
	}
}

// TestScoutThePeriphery_NextArsenalAttackReturnsBonus: when the queued attack action is itself
// played from arsenal the per-variant bonus fires (Red +3, Yellow +2, Blue +1).
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
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 0), FromArsenal: true}}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
