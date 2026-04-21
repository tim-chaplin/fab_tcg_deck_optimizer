package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfSilphidae_PlayFizzlesWithoutAura: no aura in s.Graveyard means the enter trigger
// can't banish anything and Play returns 0. AuraCreated still fires (Silphidae IS an aura).
func TestSigilOfSilphidae_PlayFizzlesWithoutAura(t *testing.T) {
	var s card.TurnState
	if got := (SigilOfSilphidaeBlue{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty graveyard)", got)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set even when banish fizzles")
	}
	if s.ArcaneDamageDealt {
		t.Errorf("ArcaneDamageDealt should stay false when banish fizzles")
	}
}

// TestSigilOfSilphidae_PlayBanishesAuraForOneArcane: an aura in s.Graveyard triggers the
// enter banish — the aura moves to Banish, Play returns 1, and ArcaneDamageDealt flips.
func TestSigilOfSilphidae_PlayBanishesAuraForOneArcane(t *testing.T) {
	aura := BlessingOfOccultRed{}
	s := card.TurnState{Graveyard: []card.Card{aura}}
	if got := (SigilOfSilphidaeBlue{}).Play(&s, nil); got != 1 {
		t.Errorf("Play() = %d, want 1", got)
	}
	if !s.ArcaneDamageDealt {
		t.Errorf("ArcaneDamageDealt should be set")
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banish)
	}
}

// TestSigilOfSilphidae_PlayNextTurnGraveyardsSelfAndFizzles: with nothing else in the
// graveyard, PlayNextTurn adds Silphidae to the graveyard and the leave trigger has no
// OTHER aura to banish — returns 0 damage.
func TestSigilOfSilphidae_PlayNextTurnGraveyardsSelfAndFizzles(t *testing.T) {
	c := SigilOfSilphidaeBlue{}
	var s card.TurnState
	r := c.PlayNextTurn(&s)
	if r.Damage != 0 {
		t.Errorf("Damage = %d, want 0 (no other aura to banish)", r.Damage)
	}
	if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != c.ID() {
		t.Errorf("Graveyard = %v, want [Silphidae]", s.Graveyard)
	}
}

// TestSigilOfSilphidae_PlayNextTurnBanishesAnotherAura: with another aura already in the
// graveyard, the leave trigger banishes it for 1 arcane. The "another" restriction is
// honoured by scan order — PlayNextTurn scans the graveyard before adding Silphidae, so the
// sigil itself can't be banished.
func TestSigilOfSilphidae_PlayNextTurnBanishesAnotherAura(t *testing.T) {
	c := SigilOfSilphidaeBlue{}
	other := BlessingOfOccultRed{}
	s := card.TurnState{Graveyard: []card.Card{other}}
	r := c.PlayNextTurn(&s)
	if r.Damage != 1 {
		t.Errorf("Damage = %d, want 1 (banished another aura)", r.Damage)
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != other.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banish)
	}
	// Silphidae stays in the graveyard after the leave trigger (it's the thing that just died).
	if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != c.ID() {
		t.Errorf("Graveyard = %v, want [Silphidae]", s.Graveyard)
	}
}
