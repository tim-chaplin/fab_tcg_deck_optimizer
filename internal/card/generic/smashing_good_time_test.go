package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSmashingGoodTime_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSmashingGoodTime_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSmashingGoodTime_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSmashingGoodTime_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (SmashingGoodTimeRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSmashingGoodTime_NextAttackReturnsBonus: first attack-action triggers +3.
func TestSmashingGoodTime_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	for _, c := range []card.Card{SmashingGoodTimeRed{}, SmashingGoodTimeYellow{}, SmashingGoodTimeBlue{}} {
		if got := c.Play(&s); got != 3 {
			t.Errorf("%s: Play() = %d, want 3", c.Name(), got)
		}
	}
}
