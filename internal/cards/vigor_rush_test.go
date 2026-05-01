package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestVigorRush_BaseGoAgainFalse pins the correctness of GoAgain() — it must be false so
// EffectiveGoAgain (baseGoAgain || self.GrantedGoAgain) can short-circuit chain-legality when
// the non-attack-action condition hasn't fired. Returning true made the chain-legality check
// always pass, over-crediting sequences where no non-attack action was played earlier in the
// turn.
func TestVigorRush_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []sim.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}} {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (gated on non-attack-action pitch)", c.Name())
		}
	}
}

// TestVigorRush_NoNonAttackActionNoGoAgain covers the miss branch: if only attack-action cards
// (or nothing) have been played this turn, the conditional go-again rider doesn't fire. The sim
// maintains NonAttackActionPlayed alongside CardsPlayed as it walks the chain, so Play reads the
// flag rather than re-scanning.
func TestVigorRush_NoNonAttackActionNoGoAgain(t *testing.T) {
	cases := []sim.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		s := sim.TurnState{
			CardsPlayed:           []sim.Card{testutils.GenericAttack(0, 0)}, // attack, not non-attack
			NonAttackActionPlayed: false,
		}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if got := s.Value; got != c.Attack() {
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
	cases := []sim.Card{VigorRushRed{}, VigorRushYellow{}, VigorRushBlue{}}
	for _, c := range cases {
		s := sim.TurnState{
			CardsPlayed:           []sim.Card{testutils.GenericAction()},
			NonAttackActionPlayed: true,
		}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (non-attack action → go again)", c.Name())
		}
	}
}
