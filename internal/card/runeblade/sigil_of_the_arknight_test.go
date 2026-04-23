package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfTheArknight_PlayOnlySetsAuraCreated verifies Play defers the reveal effect — it
// flips AuraCreated, registers a TriggerStartOfTurn entry, and returns 0. The deck peek
// happens when the sim fires the trigger next turn.
func TestSigilOfTheArknight_PlayOnlySetsAuraCreated(t *testing.T) {
	s := card.TurnState{Deck: []card.Card{stubRunebladeAttack{}}}
	if got := (SigilOfTheArknightBlue{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (reveal deferred to trigger)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
	if len(s.AuraTriggers) != 1 || s.AuraTriggers[0].Type != card.TriggerStartOfTurn {
		t.Errorf("AuraTriggers = %+v, want one TriggerStartOfTurn entry", s.AuraTriggers)
	}
}

// TestSigilOfTheArknight_TriggerRevealsAttackActionIntoHand: the post-draw deck's top card
// is an attack action → the handler moves it to s.Revealed and pops s.Deck so the deck
// loop appends the card to that turn's hand. Damage stays 0 (tempo is captured by the
// extra card, not a flat credit).
func TestSigilOfTheArknight_TriggerRevealsAttackActionIntoHand(t *testing.T) {
	var play card.TurnState
	(SigilOfTheArknightBlue{}).Play(&play, &card.CardState{})
	top := stubRunebladeAttack{}
	next := card.TurnState{Deck: []card.Card{top, stubNonAttack{}}}
	if got := play.AuraTriggers[0].Handler(&next); got != 0 {
		t.Errorf("handler damage = %d, want 0 (tempo credited via Revealed, not damage)", got)
	}
	if len(next.Revealed) != 1 || next.Revealed[0] != top {
		t.Errorf("Revealed = %v, want [%v] (top of post-draw deck)", next.Revealed, top)
	}
	if len(next.Deck) != 1 || next.Deck[0] != (stubNonAttack{}) {
		t.Errorf("Deck = %v, want top popped leaving [stubNonAttack]", next.Deck)
	}
}

// TestSigilOfTheArknight_TriggerRevealsNonAttack: top card is non-attack → Revealed stays
// nil and Deck is untouched (the card stays on top of the deck in the real game).
func TestSigilOfTheArknight_TriggerRevealsNonAttack(t *testing.T) {
	var play card.TurnState
	(SigilOfTheArknightBlue{}).Play(&play, &card.CardState{})
	next := card.TurnState{Deck: []card.Card{stubAura{}, stubRunebladeAttack{}}}
	if got := play.AuraTriggers[0].Handler(&next); got != 0 {
		t.Errorf("handler damage = %d, want 0", got)
	}
	if next.Revealed != nil {
		t.Errorf("Revealed = %v, want nil (non-attack top, no reveal)", next.Revealed)
	}
	if len(next.Deck) != 2 {
		t.Errorf("Deck len = %d, want 2 (non-attack tops aren't moved)", len(next.Deck))
	}
}

// TestSigilOfTheArknight_TriggerEmptyDeck: nothing to reveal → zero result, Revealed stays nil.
func TestSigilOfTheArknight_TriggerEmptyDeck(t *testing.T) {
	var play card.TurnState
	(SigilOfTheArknightBlue{}).Play(&play, &card.CardState{})
	var next card.TurnState
	if got := play.AuraTriggers[0].Handler(&next); got != 0 {
		t.Errorf("handler damage = %d, want 0", got)
	}
	if next.Revealed != nil {
		t.Errorf("Revealed = %v, want nil (empty deck)", next.Revealed)
	}
}

// TestSigilOfTheArknight_ImplementsAddsFutureValue pins the marker so the solver's
// beatsBest tiebreaker counts this card as future-value-adding.
func TestSigilOfTheArknight_ImplementsAddsFutureValue(t *testing.T) {
	var c card.Card = SigilOfTheArknightBlue{}
	if _, ok := c.(card.AddsFutureValue); !ok {
		t.Error("SigilOfTheArknightBlue should implement card.AddsFutureValue")
	}
}
