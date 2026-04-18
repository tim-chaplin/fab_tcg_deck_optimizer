package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestVigorRush_NoNonAttackActionNoGoAgain covers the miss branch: if only attack-action cards
// (or nothing) have been played this turn, the conditional go-again rider doesn't fire.
func TestVigorRush_NoNonAttackActionNoGoAgain(t *testing.T) {
	cases := []card.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		pc := &card.PlayedCard{Card: c}
		s := card.TurnState{
			Self:        pc,
			CardsPlayed: []card.Card{stubGenericAttack(0, 0)}, // attack action, not non-attack
		}
		if got := c.Play(&s); got != c.Attack() {
			t.Errorf("%s: Play() = %d, want %d (base power)", c.Name(), got, c.Attack())
		}
		if pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no non-attack action played)", c.Name())
		}
	}
}

// TestVigorRush_NonAttackActionGrantsGoAgain exercises the hit branch: a non-attack action played
// earlier this turn flips Self.GrantedGoAgain.
func TestVigorRush_NonAttackActionGrantsGoAgain(t *testing.T) {
	cases := []card.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		pc := &card.PlayedCard{Card: c}
		s := card.TurnState{
			Self:        pc,
			CardsPlayed: []card.Card{stubGenericAction()},
		}
		_ = c.Play(&s)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (non-attack action → go again)", c.Name())
		}
	}
}
