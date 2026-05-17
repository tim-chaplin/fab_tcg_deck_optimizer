package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// TestSigilOfSilphidae_PlayFizzlesWithoutAura: no aura in gs.Graveyard means the enter trigger
// can't banish anything and Play returns 0. AuraCreated still fires (Silphidae IS an aura)
// and a start-of-turn Aura is registered for the "destroy this" clause.
func TestSigilOfSilphidae_PlayFizzlesWithoutAura(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty graveyard)", got)
	}
	if !ge.AuraCreated() {
		t.Errorf("AuraCreated should be set even when banish fizzles")
	}
	if ge.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should stay false when banish fizzles")
	}
	if len(ge.Auras()) != 1 || ge.Auras()[0].TriggerType() != triggertype.StartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", ge.Auras())
	}
}

// TestSigilOfSilphidae_PlayBanishesAuraForOneArcane: an aura in gs.Graveyard triggers the
// enter banish — the aura moves to Banish, Play returns 1, and ArcaneDamageDealt flips.
func TestSigilOfSilphidae_PlayBanishesAuraForOneArcane(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1", got)
	}
	if !ge.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should be set")
	}
	if len(ge.Banished()) != 1 || ge.Banished()[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", ge.Banished())
	}
}

// TestSigilOfSilphidae_StartOfTurnHandlerFizzlesWithoutAnotherAura: with nothing else in the
// start-of-turn graveyard, the leave trigger has no OTHER aura to banish — handler returns
// 0 damage.
func TestSigilOfSilphidae_StartOfTurnHandlerFizzlesWithoutAnotherAura(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	next := gameengine.New()
	next.CreateAura(play.Auras()[0])
	next.FireStartOfTurn(nil)
	if next.Value() != 0 {
		t.Errorf("handler Value = %d, want 0 (no other aura to banish)", next.Value())
	}
}

// TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura: with another aura already in
// the start-of-turn graveyard, the leave trigger banishes it for 1 arcane.
func TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	other := cards.BlessingOfOccultRed{}
	next := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{other}).Build()}
	next.CreateAura(play.Auras()[0])
	next.FireStartOfTurn(nil)
	if next.Value() != 1 {
		t.Errorf("handler Value = %d, want 1 (banished another aura)", next.Value())
	}
	if len(next.Banished()) != 1 || next.Banished()[0].ID() != other.ID() {
		t.Errorf("Banish = %v, want [Blessing]", next.Banished())
	}
}
