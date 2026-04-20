package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfTheArknight_PlayOnlySetsAuraCreated verifies Play defers the reveal effect — it
// only flips AuraCreated and returns 0. The deck peek happens in PlayNextTurn now.
func TestSigilOfTheArknight_PlayOnlySetsAuraCreated(t *testing.T) {
	s := card.TurnState{Deck: []card.Card{stubRunebladeAttack{}}}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (reveal deferred to PlayNextTurn)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
}

// TestSigilOfTheArknight_PlayNextTurnRevealsAttackAction: the post-draw deck's top card is an
// attack action → credit DrawValue.
func TestSigilOfTheArknight_PlayNextTurnRevealsAttackAction(t *testing.T) {
	s := card.TurnState{Deck: []card.Card{stubRunebladeAttack{}, stubNonAttack{}}}
	if got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s); got != card.DrawValue {
		t.Errorf("PlayNextTurn() = %d, want %d (top card is an attack action)", got, card.DrawValue)
	}
}

// TestSigilOfTheArknight_PlayNextTurnRevealsNonAttack: top card is non-attack → 0.
func TestSigilOfTheArknight_PlayNextTurnRevealsNonAttack(t *testing.T) {
	s := card.TurnState{Deck: []card.Card{stubAura{}, stubRunebladeAttack{}}}
	if got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s); got != 0 {
		t.Errorf("PlayNextTurn() = %d, want 0 (top card is non-attack)", got)
	}
}

// TestSigilOfTheArknight_PlayNextTurnEmptyDeck: nothing to reveal → 0.
func TestSigilOfTheArknight_PlayNextTurnEmptyDeck(t *testing.T) {
	s := card.TurnState{}
	if got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s); got != 0 {
		t.Errorf("PlayNextTurn() = %d, want 0 (empty deck)", got)
	}
}

// TestSigilOfTheArknight_ImplementsDelayedPlay pins the marker so the deck loop queues this
// card for a next-turn callback.
func TestSigilOfTheArknight_ImplementsDelayedPlay(t *testing.T) {
	var c card.Card = SigilOfTheArknightBlue{}
	if _, ok := c.(card.DelayedPlay); !ok {
		t.Error("SigilOfTheArknightBlue should implement card.DelayedPlay")
	}
}
