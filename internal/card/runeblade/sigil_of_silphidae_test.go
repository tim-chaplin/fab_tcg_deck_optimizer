package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfSilphidae_PlayFizzlesWithoutAura: no aura in s.Graveyard means the enter trigger
// can't banish anything and Play returns 0. AuraCreated still fires (Silphidae IS an aura)
// and a start-of-turn AuraTrigger is registered for the "destroy this" clause.
func TestSigilOfSilphidae_PlayFizzlesWithoutAura(t *testing.T) {
	var s card.TurnState
	(SigilOfSilphidaeBlue{}).Play(&s, &card.CardState{Card: SigilOfSilphidaeBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (empty graveyard)", got)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set even when banish fizzles")
	}
	if s.ArcaneDamageDealt {
		t.Errorf("ArcaneDamageDealt should stay false when banish fizzles")
	}
	if len(s.AuraTriggers) != 1 || s.AuraTriggers[0].Type != card.TriggerStartOfTurn {
		t.Errorf("AuraTriggers = %+v, want one TriggerStartOfTurn entry", s.AuraTriggers)
	}
}

// TestSigilOfSilphidae_PlayBanishesAuraForOneArcane: an aura in s.Graveyard triggers the
// enter banish — the aura moves to Banish, Play returns 1, and ArcaneDamageDealt flips.
func TestSigilOfSilphidae_PlayBanishesAuraForOneArcane(t *testing.T) {
	aura := BlessingOfOccultRed{}
	var s card.TurnState
	s.SetGraveyard([]card.Card{aura})
	(SigilOfSilphidaeBlue{}).Play(&s, &card.CardState{Card: SigilOfSilphidaeBlue{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1", got)
	}
	if !s.ArcaneDamageDealt {
		t.Errorf("ArcaneDamageDealt should be set")
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banish)
	}
}

// TestSigilOfSilphidae_StartOfTurnHandlerFizzlesWithoutAnotherAura: with nothing else in the
// start-of-turn graveyard, the leave trigger has no OTHER aura to banish — handler returns
// 0 damage.
func TestSigilOfSilphidae_StartOfTurnHandlerFizzlesWithoutAnotherAura(t *testing.T) {
	var play card.TurnState
	(SigilOfSilphidaeBlue{}).Play(&play, &card.CardState{Card: SigilOfSilphidaeBlue{}})
	var next card.TurnState
	got := play.AuraTriggers[0].Handler(&next)
	if got != 0 {
		t.Errorf("handler damage = %d, want 0 (no other aura to banish)", got)
	}
}

// TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura: with another aura already in
// the start-of-turn graveyard, the leave trigger banishes it for 1 arcane. The sim
// graveyards Self only AFTER this handler returns, so the scan can't pick up Silphidae
// itself — the printed "another aura" restriction is satisfied naturally.
func TestSigilOfSilphidae_StartOfTurnHandlerBanishesAnotherAura(t *testing.T) {
	var play card.TurnState
	(SigilOfSilphidaeBlue{}).Play(&play, &card.CardState{Card: SigilOfSilphidaeBlue{}})
	other := BlessingOfOccultRed{}
	var next card.TurnState
	next.SetGraveyard([]card.Card{other})
	got := play.AuraTriggers[0].Handler(&next)
	if got != 1 {
		t.Errorf("handler damage = %d, want 1 (banished another aura)", got)
	}
	if len(next.Banish) != 1 || next.Banish[0].ID() != other.ID() {
		t.Errorf("Banish = %v, want [Blessing]", next.Banish)
	}
}
