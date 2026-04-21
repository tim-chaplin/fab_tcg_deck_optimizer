package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestComeToFight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestComeToFight_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{ComeToFightRed{}, ComeToFightYellow{}, ComeToFightBlue{}} {
		if got := c.Play(&s, nil); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestComeToFight_NonAttackInRemainingFizzles: only an action (no attack) in CardsRemaining — the
// attack-action predicate rejects it.
func TestComeToFight_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (ComeToFightRed{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestComeToFight_NextAttackReturnsBonus: first attack-action in CardsRemaining triggers the
// per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestComeToFight_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	cases := []struct {
		c    card.Card
		want int
	}{
		{ComeToFightRed{}, 3},
		{ComeToFightYellow{}, 2},
		{ComeToFightBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
