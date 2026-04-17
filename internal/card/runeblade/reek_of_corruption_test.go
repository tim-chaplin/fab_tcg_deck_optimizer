package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestReekOfCorruption_NoAuraReturnsBaseAttack covers the fix: without an aura created or played
// this turn there's no satisfying the "played or created an aura" clause, so the discard rider
// doesn't fire and Play returns the printed attack alone.
func TestReekOfCorruption_NoAuraReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ReekOfCorruptionRed{}, 4},
		{ReekOfCorruptionYellow{}, 3},
		{ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, no aura)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_AuraCreatedTriggersDiscard exercises the AuraCreated branch: a prior card
// (e.g. Reduce to Runechant) setting AuraCreated satisfies the rider even if no aura card is in
// CardsPlayed.
func TestReekOfCorruption_AuraCreatedTriggersDiscard(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ReekOfCorruptionRed{}, 4 + 3},
		{ReekOfCorruptionYellow{}, 3 + 3},
		{ReekOfCorruptionBlue{}, 2 + 3},
	}
	for _, tc := range cases {
		s := card.TurnState{AuraCreated: true}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + discard rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_AuraPlayedTriggersDiscard exercises the HasPlayedType branch: an aura-
// typed card earlier in the chain also satisfies "played an aura this turn".
func TestReekOfCorruption_AuraPlayedTriggersDiscard(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (ReekOfCorruptionRed{}).Play(&s); got != 4+3 {
		t.Errorf("Play() = %d, want %d (attack + discard rider, aura earlier in chain)", got, 4+3)
	}
}
