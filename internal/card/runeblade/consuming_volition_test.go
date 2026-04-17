package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestConsumingVolition_NoRunechantsReturnsBaseAttack covers the fix: if no Runechant exists at
// Play time there's no source of arcane damage this turn, so the discard rider can't fire and
// Play returns only the printed attack — no flat +3.
func TestConsumingVolition_NoRunechantsReturnsBaseAttack(t *testing.T) {
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
			t.Errorf("%s: Play() = %d, want %d (base attack, no runechants)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_RunechantsTriggerDiscardRider exercises the satisfied path: any live
// Runechant means arcane damage will fire on this attack, so the discard rider activates and
// Play returns attack + 3.
func TestConsumingVolition_RunechantsTriggerDiscardRider(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ConsumingVolitionRed{}, 4 + 3},
		{ConsumingVolitionYellow{}, 3 + 3},
		{ConsumingVolitionBlue{}, 2 + 3},
	}
	for _, tc := range cases {
		s := card.TurnState{Runechants: 1}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + discard rider)", tc.c.Name(), got, tc.want)
		}
	}
}
