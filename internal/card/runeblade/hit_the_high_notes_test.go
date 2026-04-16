package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestHitTheHighNotes_NoAuraReturnsBase(t *testing.T) {
	// Neither an aura played nor one created this turn → no bonus, just printed power.
	cases := []struct {
		c    card.Card
		base int
	}{
		{HitTheHighNotesRed{}, 4},
		{HitTheHighNotesYellow{}, 3},
		{HitTheHighNotesBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
	}
}

func TestHitTheHighNotes_AuraPlayedTriggersBonus(t *testing.T) {
	// An Aura-typed card earlier in the turn's CardsPlayed → +2 power.
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (HitTheHighNotesRed{}).Play(&s); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 aura bonus)", got)
	}
}

func TestHitTheHighNotes_AuraCreatedTriggersBonus(t *testing.T) {
	// AuraCreated flag set earlier in the chain (e.g. Runechant creation) → +2 power, even
	// without an Aura-typed card in CardsPlayed.
	s := card.TurnState{AuraCreated: true}
	if got := (HitTheHighNotesRed{}).Play(&s); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 AuraCreated bonus)", got)
	}
}
