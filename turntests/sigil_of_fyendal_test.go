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

// TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger: Play flips AuraCreated and appends a
// start-of-turn Aura with Count=1 — no same-turn damage, the 1{h} gain is credited
// when the sim fires the trigger next turn.
func TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SigilOfFyendalBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (1{h} gain deferred to trigger)", got)
	}
	if !ge.AuraCreated() {
		t.Error("AuraCreated = false, want true")
	}
	if len(ge.Auras()) != 1 || ge.Auras()[0].TriggerType() != triggertype.StartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", ge.Auras())
	}
	if ge.Auras()[0].Count() != 1 {
		t.Errorf("Count = %d, want 1", ge.Auras()[0].Count())
	}
}

// Tests that a carried Sigil of Fyendal fires at the start of the next turn, crediting the
// 1{h} gain (worth 1 damage) and self-destructing.
func TestSigilOfFyendal_StartOfTurnCredits1Damage(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfFyendalBlue{}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 1 {
		t.Errorf("Value = %d, want 1 (start-of-turn 1{h} gain)", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Sigil of Fyendal self-destructed after firing)", got)
	}
}
