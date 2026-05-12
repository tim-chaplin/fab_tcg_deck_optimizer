package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestSigilOfSilphidae_PlayFizzlesWithoutAura: no aura in s.Graveyard means the enter trigger
// can't banish anything and Play returns 0. AuraCreated still fires (Silphidae IS an aura)
// and a start-of-turn Aura is registered for the "destroy this" clause.
func TestSigilOfSilphidae_PlayFizzlesWithoutAura(t *testing.T) {
	var s sim.TurnState
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty graveyard)", got)
	}
	if !s.AuraCreated() {
		t.Errorf("AuraCreated should be set even when banish fizzles")
	}
	if s.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should stay false when banish fizzles")
	}
	if len(s.Auras()) != 1 || s.Auras()[0].TriggerType != sim.TriggerStartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", s.Auras())
	}
}

// TestSigilOfSilphidae_PlayBanishesAuraForOneArcane: an aura in s.Graveyard triggers the
// enter banish — the aura moves to Banish, Play returns 1, and ArcaneDamageDealt flips.
func TestSigilOfSilphidae_PlayBanishesAuraForOneArcane(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	s := sim.NewTurnStateFromCards(nil, []card.Card{aura})
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
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
	var play sim.TurnState
	sim.ResolveChainStep(&play, play.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	next := sim.NewTurnStateFromCards(nil, nil)
	next.SetAuras(append(next.Auras(), play.Auras()[0]))
	next.SetCurrentAuraIdxForTesting(0)
	next.FireAuraForTesting(0)
	if next.Value() != 0 {
		t.Errorf("handler Value = %d, want 0 (no other aura to banish)", next.Value())
	}
}

// TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura: with another aura already in
// the start-of-turn graveyard, the leave trigger banishes it for 1 arcane.
func TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura(t *testing.T) {
	var play sim.TurnState
	sim.ResolveChainStep(&play, play.Logger(), &card.CardState{Card: cards.SigilOfSilphidaeBlue{}})
	other := cards.BlessingOfOccultRed{}
	next := sim.NewTurnStateFromCards(nil, []card.Card{other})
	next.SetAuras(append(next.Auras(), play.Auras()[0]))
	next.SetCurrentAuraIdxForTesting(0)
	next.FireAuraForTesting(0)
	if next.Value() != 1 {
		t.Errorf("handler Value = %d, want 1 (banished another aura)", next.Value())
	}
	if len(next.Banished()) != 1 || next.Banished()[0].ID() != other.ID() {
		t.Errorf("Banish = %v, want [Blessing]", next.Banished())
	}
}
