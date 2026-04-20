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

// TestSigilOfTheArknight_PlayNextTurnRevealsAttackActionIntoHand: the post-draw deck's top card
// is an attack action → the card is returned in ToHand so the deck loop moves it into the hand.
// Damage stays 0 because the tempo is captured by having the extra card, not by a flat credit.
func TestSigilOfTheArknight_PlayNextTurnRevealsAttackActionIntoHand(t *testing.T) {
	top := stubRunebladeAttack{}
	s := card.TurnState{Deck: []card.Card{top, stubNonAttack{}}}
	got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s)
	if got.Damage != 0 {
		t.Errorf("Damage = %d, want 0 (tempo credited via ToHand, not Damage)", got.Damage)
	}
	if got.ToHand != top {
		t.Errorf("ToHand = %v, want %v (top of post-draw deck)", got.ToHand, top)
	}
}

// TestSigilOfTheArknight_PlayNextTurnRevealsNonAttack: top card is non-attack → ToHand stays
// nil (the card stays on top of the deck in the real game).
func TestSigilOfTheArknight_PlayNextTurnRevealsNonAttack(t *testing.T) {
	s := card.TurnState{Deck: []card.Card{stubAura{}, stubRunebladeAttack{}}}
	got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s)
	if got.ToHand != nil {
		t.Errorf("ToHand = %v, want nil (top is non-attack, no reveal)", got.ToHand)
	}
	if got.Damage != 0 {
		t.Errorf("Damage = %d, want 0", got.Damage)
	}
}

// TestSigilOfTheArknight_PlayNextTurnEmptyDeck: nothing to reveal → zero result.
func TestSigilOfTheArknight_PlayNextTurnEmptyDeck(t *testing.T) {
	s := card.TurnState{}
	got := (SigilOfTheArknightBlue{}).PlayNextTurn(&s)
	if got.ToHand != nil || got.Damage != 0 {
		t.Errorf("PlayNextTurn() = %+v, want zero (empty deck)", got)
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
