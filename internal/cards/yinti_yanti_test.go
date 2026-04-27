package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestYintiYanti_NoAuraReturnsBase covers the miss branch: without an aura played or created this
// turn the attacking-side +1{p} rider doesn't fire and Play returns the printed power.
func TestYintiYanti_NoAuraReturnsBase(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{YintiYantiRed{}, 3},
		{YintiYantiYellow{}, 2},
		{YintiYantiBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (no aura → base power)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestYintiYanti_AuraCreatedAddsOne exercises the AuraCreated branch: an aura created earlier in
// the same chain (sets the flag on TurnState) triggers the +1{p} rider.
func TestYintiYanti_AuraCreatedAddsOne(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{YintiYantiRed{}, 4},
		{YintiYantiYellow{}, 3},
		{YintiYantiBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{AuraCreated: true}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (AuraCreated → +1)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestYintiYanti_AuraPlayedAddsOne exercises the HasPlayedType branch: an aura-typed card earlier
// in CardsPlayed also satisfies the rider even if AuraCreated is still false (e.g. an Aura card
// that entered via a non-creation path).
func TestYintiYanti_AuraPlayedAddsOne(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubGenericAura()}}
	(YintiYantiRed{}).Play(&s, &card.CardState{Card: YintiYantiRed{}})
	if got := s.Value; got != 4 {
		t.Errorf("Play() = %d, want 4 (aura in CardsPlayed → +1)", got)
	}
}
