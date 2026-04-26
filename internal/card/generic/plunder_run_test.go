package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestPlunderRun_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestPlunderRun_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{PlunderRunRed{}, PlunderRunYellow{}, PlunderRunBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPlunderRun_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPlunderRun_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	(PlunderRunRed{}).Play(&s, &card.CardState{Card: PlunderRunRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPlunderRun_NextAttackGrantsBonusAttack: arsenal-played copy with a queued attack
// action grants the per-variant bonus (Red +3, Yellow +2, Blue +1) onto the target's
// BonusAttack — granter returns 0; the +N attributes to the buffed attack.
func TestPlunderRun_NextAttackGrantsBonusAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{PlunderRunRed{}, 3},
		{PlunderRunYellow{}, 2},
		{PlunderRunBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(0, 0)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		self := &card.CardState{Card: tc.c, FromArsenal: true}
		tc.c.Play(&s, self)
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// TestPlunderRun_HandPlayedFizzles: hand-played copy fails the from-arsenal gate even when a
// queued attack action would otherwise satisfy the rider.
func TestPlunderRun_HandPlayedFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 0)}}}
	self := &card.CardState{Card: PlunderRunRed{}}
	(PlunderRunRed{}).Play(&s, self)
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (hand-played, not from arsenal)", got)
	}
}
