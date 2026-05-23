package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
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

// Tests that a carried Sigil of the Arknight whose post-self-destruct deck top is an
// attack action draws it into the hand — the drawn attack then plays this turn, observable
// as its damage contribution to Value (the hand otherwise carries no attack).
func TestSigilOfTheArknight_StartOfTurnRevealsAttackActionIntoHand(t *testing.T) {
	revealed := testutils.AttackWithPower{Power: 7}
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfTheArknightBlue{}).
		Build()
	deckCards := []deck.Card{revealed}
	deckCards = append(deckCards, nil...)
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 7 {
		t.Errorf("Value = %d, want 7 (revealed attack landed for 7)", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Sigil self-destructed)", got)
	}
}

// Tests that a carried Sigil of the Arknight whose post-self-destruct deck top is a
// non-attack-action card draws nothing — the turn's value stays at zero because no attack
// was revealed and the hand carries none.
func TestSigilOfTheArknight_StartOfTurnRevealsNonAttackDoesNotDraw(t *testing.T) {
	top := testutils.NonAttack{}
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfTheArknightBlue{}).
		Build()
	deckCards := []deck.Card{top}
	deckCards = append(deckCards, nil...)
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 0 {
		t.Errorf("Value = %d, want 0 (non-attack top isn't drawn)", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Sigil self-destructs even when reveal whiffs)", got)
	}
}

// Tests that a carried Sigil of the Arknight with an empty deck still self-destructs at
// start of turn — the reveal is a silent no-op.
func TestSigilOfTheArknight_StartOfTurnEmptyDeckSelfDestructs(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfTheArknightBlue{}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := len(summary.State.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Sigil self-destructs even with empty deck)", got)
	}
}
