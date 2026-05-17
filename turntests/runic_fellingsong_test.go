package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly: an empty graveyard fizzles the banish
// rider, so Play returns just Attack().
func TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly(t *testing.T) {
	ge := gameengine.New()
	c := cards.RunicFellingsongRed{}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
	if got := ge.Value(); got != c.Attack() {
		t.Errorf("Play() = %d, want %d (Attack only; banish fizzles)", got, c.Attack())
	}
}

// TestRunicFellingsong_AuraInGraveyardFiresBanishRider: with an aura banishable, Play returns
// Attack() + 1 (the banish rider's arcane).
func TestRunicFellingsong_AuraInGraveyardFiresBanishRider(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura}).Build()}
	c := cards.RunicFellingsongRed{}
	want := c.Attack() + 1
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
	if got := ge.Value(); got != want {
		t.Errorf("Play() = %d, want %d (Attack + banish rider)", got, want)
	}
	if len(ge.Banished()) != 1 || ge.Banished()[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", ge.Banished())
	}
}
