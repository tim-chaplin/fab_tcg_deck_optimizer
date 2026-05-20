package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
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

// TestSigilOfFyendal_TriggerHandlerCredits1Damage: the registered handler credits +1 damage
// (the 1{h} gain, valued 1-to-1 with damage).
func TestSigilOfFyendal_TriggerHandlerCredits1Damage(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SigilOfFyendalBlue{}})
	fire := gameengine.New()
	fire.CreateAura(ge.Auras()[0])
	fire.FireTriggers(triggertype.StartOfTurn, nil)
	if fire.Value() != 1 {
		t.Errorf("Handler Value = %d, want 1", fire.Value())
	}
}
