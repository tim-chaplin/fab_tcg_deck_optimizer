package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestScoutThePeriphery_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestScoutThePeriphery_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{ScoutThePeripheryRed{}, ScoutThePeripheryYellow{}, ScoutThePeripheryBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestScoutThePeriphery_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestScoutThePeriphery_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (ScoutThePeripheryRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestScoutThePeriphery_NextAttackReturnsBonus: first attack-action triggers +3.
func TestScoutThePeriphery_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	for _, c := range []card.Card{ScoutThePeripheryRed{}, ScoutThePeripheryYellow{}, ScoutThePeripheryBlue{}} {
		if got := c.Play(&s); got != 3 {
			t.Errorf("%s: Play() = %d, want 3", c.Name(), got)
		}
	}
}
