package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly: an empty graveyard fizzles the banish
// rider, so summary.Value is just Attack().
func TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly(t *testing.T) {
	c := cards.RunicFellingsongRed{}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{c, testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != c.Attack() {
		t.Errorf("Value = %d, want %d (Attack only; banish fizzles)", summary.Value, c.Attack())
	}
}

// TestRunicFellingsong_AuraInGraveyardFiresBanishRider: with an aura banishable, Play credits
// Attack() + 1 (the banish rider's arcane) and the aura ends up in the banish zone.
func TestRunicFellingsong_AuraInGraveyardFiresBanishRider(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	c := cards.RunicFellingsongRed{}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	prior := gameengine.GameStateBuilder().
		SetGraveyard([]card.Card{aura}).
		SetIncomingDamage(0).
		Build()
	hand := []card.Card{c, testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, prior, hand)
	if want := c.Attack() + 1; summary.Value != want {
		t.Errorf("Value = %d, want %d (Attack + banish rider)", summary.Value, want)
	}
	if b := summary.State.Banished(); len(b) != 1 || b[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", b)
	}
}
