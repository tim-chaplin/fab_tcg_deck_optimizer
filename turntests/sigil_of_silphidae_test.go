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

// Tests that a carried Sigil of Silphidae with an empty graveyard self-destructs at start
// of turn but its leave trigger fizzles — no arcane damage, no banish.
func TestSigilOfSilphidae_StartOfTurnFizzlesWithoutAnotherAura(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfSilphidaeBlue{}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 0 {
		t.Errorf("Value = %d, want 0 (leave trigger fizzles with empty graveyard)", got)
	}
	if got := len(summary.State.Banished()); got != 0 {
		t.Errorf("Banished = %d, want 0", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Sigil self-destructed)", got)
	}
}

// Tests that a carried Sigil of Silphidae with a banishable aura in the graveyard banishes
// it for 1 arcane damage when it self-destructs at start of turn.
func TestSigilOfSilphidae_StartOfTurnBanishesAnotherAura(t *testing.T) {
	other := cards.BlessingOfOccultRed{}
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfSilphidaeBlue{}).
		SetGraveyard([]card.Card{other}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 1 {
		t.Errorf("Value = %d, want 1 (leave trigger banished another aura for 1 arcane)", got)
	}
	if b := summary.State.Banished(); len(b) != 1 || b[0].ID() != other.ID() {
		t.Errorf("Banished = %v, want [Blessing]", b)
	}
}
