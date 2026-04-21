package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRunicFellingsong_NoAuraCreditsPrintedArcaneOnly: an empty graveyard fizzles the banish
// rider. Play returns Attack() + 1 for the printed arcane.
func TestRunicFellingsong_NoAuraCreditsPrintedArcaneOnly(t *testing.T) {
	var s card.TurnState
	c := RunicFellingsongRed{}
	want := c.Attack() + 1
	if got := c.Play(&s, nil); got != want {
		t.Errorf("Play() = %d, want %d (Attack + 1 printed arcane)", got, want)
	}
}

// TestRunicFellingsong_AuraInGraveyardFiresBanishRider: with an aura banishable, Play returns
// Attack() + 2 (printed arcane + banish rider's 1 arcane).
func TestRunicFellingsong_AuraInGraveyardFiresBanishRider(t *testing.T) {
	aura := BlessingOfOccultRed{}
	s := card.TurnState{Graveyard: []card.Card{aura}}
	c := RunicFellingsongRed{}
	want := c.Attack() + 2
	if got := c.Play(&s, nil); got != want {
		t.Errorf("Play() = %d, want %d (Attack + printed arcane + banish rider)", got, want)
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != aura.ID() {
		t.Errorf("Banish = %v, want [Blessing]", s.Banish)
	}
}
