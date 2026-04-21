package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSmashingGoodTime_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSmashingGoodTime_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSmashingGoodTime_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSmashingGoodTime_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	if got := (SmashingGoodTimeRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSmashingGoodTime_NextAttackReturnsBonus: arsenal-played copy with a queued attack action
// triggers the per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestSmashingGoodTime_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SmashingGoodTimeRed{}, 3},
		{SmashingGoodTimeYellow{}, 2},
		{SmashingGoodTimeBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 0)}}}
		self := &card.CardState{Card: tc.c, FromArsenal: true}
		if got := tc.c.Play(&s, self); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

// TestSmashingGoodTime_HandPlayedFizzles: hand-played copy fails the from-arsenal gate.
func TestSmashingGoodTime_HandPlayedFizzles(t *testing.T) {
	for _, c := range []card.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 0)}}}
		self := &card.CardState{Card: c}
		if got := c.Play(&s, self); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (hand-played)", c.Name(), got)
		}
	}
}
