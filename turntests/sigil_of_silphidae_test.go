package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// TestSigilOfSilphidae_PlayFizzlesWithoutAura: no aura in s.Graveyard means the enter trigger
// can't banish anything and Play returns 0. AuraCreated still fires (Silphidae IS an aura)
// and a start-of-turn Aura is registered for the "destroy this" clause.
func TestSigilOfSilphidae_PlayFizzlesWithoutAura(t *testing.T) {
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty graveyard)", got)
	}
	if !s.AuraCreated() {
		t.Errorf("AuraCreated should be set even when banish fizzles")
	}
	if s.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should stay false when banish fizzles")
	}
	if len(s.Auras()) != 1 || s.Auras()[0].TriggerType() != triggertype.StartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", s.Auras())
	}
}

// TestSigilOfSilphidae_PlayBanishesAuraForOneArcane: an aura in s.Graveyard triggers the
// enter banish — the aura moves to Banish, Play returns 1, and ArcaneDamageDealt flips.
func TestSigilOfSilphidae_PlayBanishesAuraForOneArcane(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1", got)
	}
	if !s.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should be set")
	}
	if len(s.Banished()) != 1 || s.Banished()[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banished())
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
