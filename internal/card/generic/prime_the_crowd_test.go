package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestPrimeTheCrowd_NoAttackReturnsZero: no qualifying next attack card → +4 rider fizzles.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{PrimeTheCrowdRed{}, PrimeTheCrowdYellow{}, PrimeTheCrowdBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPrimeTheCrowd_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (PrimeTheCrowdRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPrimeTheCrowd_NextAttackReturnsBonus: first attack-action triggers +4.
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	for _, c := range []card.Card{PrimeTheCrowdRed{}, PrimeTheCrowdYellow{}, PrimeTheCrowdBlue{}} {
		if got := c.Play(&s); got != 4 {
			t.Errorf("%s: Play() = %d, want 4", c.Name(), got)
		}
	}
}
