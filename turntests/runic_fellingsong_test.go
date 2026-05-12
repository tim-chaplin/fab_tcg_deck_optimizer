package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly: an empty graveyard fizzles the banish
// rider, so Play returns just Attack().
func TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly(t *testing.T) {
	s := gameengine.New()
	c := cards.RunicFellingsongRed{}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
	if got := s.Value(); got != c.Attack() {
		t.Errorf("Play() = %d, want %d (Attack only; banish fizzles)", got, c.Attack())
	}
}

// TestRunicFellingsong_AuraInGraveyardFiresBanishRider: with an aura banishable, Play returns
// Attack() + 1 (the banish rider's arcane).
func TestRunicFellingsong_AuraInGraveyardFiresBanishRider(t *testing.T) {
	aura := cards.BlessingOfOccultRed{}
	s := gameengine.NewFromCards(nil, []card.Card{aura})
	c := cards.RunicFellingsongRed{}
	want := c.Attack() + 1
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
	if got := s.Value(); got != want {
		t.Errorf("Play() = %d, want %d (Attack + banish rider)", got, want)
	}
	if len(s.Banished()) != 1 || s.Banished()[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banished())
	}
}
