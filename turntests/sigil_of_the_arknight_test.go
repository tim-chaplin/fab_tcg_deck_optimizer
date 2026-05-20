package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// TestSigilOfTheArknight_PlayOnlySetsAuraCreated verifies Play defers the reveal effect — it
// flips AuraCreated, registers a TriggerStartOfTurn entry, and returns 0. The deck peek
// happens when the sim fires the trigger next turn.
func TestSigilOfTheArknight_PlayOnlySetsAuraCreated(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.RunebladeAttack{}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (reveal deferred to trigger)", got)
	}
	if !ge.AuraCreated() {
		t.Error("AuraCreated = false, want true")
	}
	if len(ge.Auras()) != 1 || ge.Auras()[0].TriggerType() != triggertype.StartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", ge.Auras())
	}
}

// TestSigilOfTheArknight_TriggerRevealsAttackActionIntoHand: the post-draw deck's top card
// is an attack action → the handler draws it into the hand and pops the deck. Damage
// stays 0 (tempo is captured by the extra card, not a flat credit).
func TestSigilOfTheArknight_TriggerRevealsAttackActionIntoHand(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	top := testutils.RunebladeAttack{}
	next := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top, testutils.NonAttack{}}).Build()}
	next.CreateAura(play.Auras()[0])
	next.FireTriggers(triggertype.StartOfTurn, nil)
	if next.Value() != 0 {
		t.Errorf("handler Value = %d, want 0 (tempo credited via the draw, not damage)", next.Value())
	}
	if h := next.Hand(); len(h) != 1 || h[0] != top {
		t.Errorf("Hand = %v, want [%v] (top of post-draw deck)", h, top)
	}
	if d := next.Deck(); d.Size() != 1 || d.PeekTop().(card.Card) != (testutils.NonAttack{}) {
		t.Errorf("Deck = %v, want top popped leaving [testutils.NonAttack]", d)
	}
}

// TestSigilOfTheArknight_TriggerRevealsNonAttack: top card is non-attack → Hand stays
// empty and Deck is untouched (the card stays on top of the deck in the real game).
func TestSigilOfTheArknight_TriggerRevealsNonAttack(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	next := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.Aura{}, testutils.RunebladeAttack{}}).Build()}
	next.CreateAura(play.Auras()[0])
	next.FireTriggers(triggertype.StartOfTurn, nil)
	if next.Value() != 0 {
		t.Errorf("handler Value = %d, want 0", next.Value())
	}
	if h := next.Hand(); h != nil {
		t.Errorf("Hand = %v, want nil (non-attack top, no draw)", h)
	}
	if d := next.Deck(); d.Size() != 2 {
		t.Errorf("Deck size = %d, want 2 (non-attack tops aren't moved)", d.Size())
	}
}

// TestSigilOfTheArknight_TriggerEmptyDeck: nothing to reveal → zero result, Hand stays empty.
func TestSigilOfTheArknight_TriggerEmptyDeck(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	next := gameengine.New()
	next.CreateAura(play.Auras()[0])
	next.FireTriggers(triggertype.StartOfTurn, nil)
	if next.Value() != 0 {
		t.Errorf("handler Value = %d, want 0", next.Value())
	}
	if h := next.Hand(); h != nil {
		t.Errorf("Hand = %v, want nil (empty deck)", h)
	}
}
