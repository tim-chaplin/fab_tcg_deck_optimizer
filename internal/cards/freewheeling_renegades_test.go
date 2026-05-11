package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that each variant credits printed_power - 2 to s.Value().
func TestFreewheelingRenegades_AlwaysDebuffedByTwo(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{FreewheelingRenegadesRed{}, 4},
		{FreewheelingRenegadesYellow{}, 3},
		{FreewheelingRenegadesBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (printed - 2)", tc.c.Name(), got, tc.want)
		}
	}
}
