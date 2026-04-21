package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestVigorRush_BaseGoAgainFalse pins the correctness of GoAgain() — it must be false so
// EffectiveGoAgain (baseGoAgain || self.GrantedGoAgain) can short-circuit chain-legality when
// the non-attack-action condition hasn't fired. Returning true made the chain-legality check
// always pass, over-crediting sequences where no non-attack action was played earlier in the
// turn.
func TestVigorRush_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []card.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}} {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (gated on non-attack-action pitch)", c.Name())
		}
	}
}

// TestVigorRush_NoNonAttackActionNoGoAgain covers the miss branch: if only attack-action cards
// (or nothing) have been played this turn, the conditional go-again rider doesn't fire.
func TestVigorRush_NoNonAttackActionNoGoAgain(t *testing.T) {
	cases := []card.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		s := card.TurnState{
			CardsPlayed: []card.Card{stubGenericAttack(0, 0)}, // attack action, not non-attack
		}
		self := &card.CardState{Card: c}
		if got := c.Play(&s, self); got != c.Attack() {
			t.Errorf("%s: Play() = %d, want %d (base power)", c.Name(), got, c.Attack())
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no non-attack action played)", c.Name())
		}
	}
}

// TestVigorRush_NonAttackActionGrantsGoAgain exercises the hit branch: a non-attack action played
// earlier this turn flips self.GrantedGoAgain.
func TestVigorRush_NonAttackActionGrantsGoAgain(t *testing.T) {
	cases := []card.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		s := card.TurnState{CardsPlayed: []card.Card{stubGenericAction()}}
		self := &card.CardState{Card: c}
		_ = c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (non-attack action → go again)", c.Name())
		}
	}
}
