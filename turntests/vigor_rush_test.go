package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestVigorRush_BaseGoAgainFalse pins GoAgain() = false so EffectiveGoAgain short-circuits
// chain-legality when the non-attack-action condition hasn't fired.
func TestVigorRush_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []card.Card{cards.VigorRushRed{}, cards.VigorRushYellow{}, cards.VigorRushBlue{}} {
		if c.GoAgain(nil) {
			t.Errorf("%s: GoAgain() = true, want false (gated on non-attack-action pitch)", c.Name())
		}
	}
}

// TestVigorRush_NoNonAttackActionNoGoAgain covers the miss branch: with only attack-action
// cards (or nothing) played this turn, the conditional go-again rider doesn't fire.
func TestVigorRush_NoNonAttackActionNoGoAgain(t *testing.T) {
	cases := []card.Card{cards.VigorRushRed{}, cards.VigorRushYellow{}, cards.VigorRushBlue{}}
	for _, c := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetCardsPlayed([]card.Card{testutils.FakeRedAttack()}). // not non-attack
			SetNonAttackActionPlayed(false).
			Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if got := ge.Value(); got != c.Attack() {
			t.Errorf("%s: Play() = %d, want %d (base power)", c.Name(), got, c.Attack())
		}
		if pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no non-attack action played)", c.Name())
		}
	}
}

// TestVigorRush_NonAttackActionGrantsGoAgain exercises the hit branch: a non-attack action played
// earlier this turn flips pc.GrantedGoAgain.
func TestVigorRush_NonAttackActionGrantsGoAgain(t *testing.T) {
	cases := []card.Card{cards.VigorRushRed{}, cards.VigorRushYellow{}, cards.VigorRushBlue{}}
	for _, c := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetCardsPlayed([]card.Card{testutils.FakeRedAction()}).
			SetNonAttackActionPlayed(true).
			Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (non-attack action → go again)", c.Name())
		}
	}
}
