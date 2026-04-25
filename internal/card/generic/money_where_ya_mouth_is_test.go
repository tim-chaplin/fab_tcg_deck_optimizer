package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMoneyWhereYaMouthIs_NoAttackReturnsZero: no qualifying next attack card → +3 fizzles.
func TestMoneyWhereYaMouthIs_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{MoneyWhereYaMouthIsRed{}, MoneyWhereYaMouthIsYellow{}, MoneyWhereYaMouthIsBlue{}} {
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestMoneyWhereYaMouthIs_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestMoneyWhereYaMouthIs_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	if got := (MoneyWhereYaMouthIsRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestMoneyWhereYaMouthIs_NextAttackReturnsBonus: first attack-action triggers the per-variant
// bonus (Red +3, Yellow +2, Blue +1).
func TestMoneyWhereYaMouthIs_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MoneyWhereYaMouthIsRed{}, 3},
		{MoneyWhereYaMouthIsYellow{}, 2},
		{MoneyWhereYaMouthIsBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(0, 0)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		if got := tc.c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
