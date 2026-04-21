package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly: an empty graveyard fizzles the banish
// rider, so Play returns just Attack().
func TestRunicFellingsong_NoAuraCreditsPrintedPowerOnly(t *testing.T) {
	var s card.TurnState
	c := RunicFellingsongRed{}
	if got := c.Play(&s, nil); got != c.Attack() {
		t.Errorf("Play() = %d, want %d (Attack only; banish fizzles)", got, c.Attack())
	}
}

// TestRunicFellingsong_AuraInGraveyardFiresBanishRider: with an aura banishable, Play returns
// Attack() + 1 (the banish rider's arcane).
func TestRunicFellingsong_AuraInGraveyardFiresBanishRider(t *testing.T) {
	aura := BlessingOfOccultRed{}
	s := card.TurnState{Graveyard: []card.Card{aura}}
	c := RunicFellingsongRed{}
	want := c.Attack() + 1
	if got := c.Play(&s, nil); got != want {
		t.Errorf("Play() = %d, want %d (Attack + banish rider)", got, want)
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banish)
	}
}
