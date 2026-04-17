package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestConsumingVolition_ArcaneDamageNotDealtReturnsBaseAttack covers the unsatisfied branch:
// when ArcaneDamageDealt is false the discard rider can't fire and Play returns only the
// printed attack.
func TestConsumingVolition_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ConsumingVolitionRed{}, 4},
		{ConsumingVolitionYellow{}, 3},
		{ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_ArcaneDamageDealtTriggersDiscardRider exercises the satisfied path:
// when ArcaneDamageDealt is set (a prior attack fired a Runechant, or a direct-arcane card
// flipped the flag), the discard rider activates and Play returns attack + 3.
func TestConsumingVolition_ArcaneDamageDealtTriggersDiscardRider(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ConsumingVolitionRed{}, 4 + 3},
		{ConsumingVolitionYellow{}, 3 + 3},
		{ConsumingVolitionBlue{}, 2 + 3},
	}
	for _, tc := range cases {
		s := card.TurnState{ArcaneDamageDealt: true}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + discard rider)", tc.c.Name(), got, tc.want)
		}
	}
}
