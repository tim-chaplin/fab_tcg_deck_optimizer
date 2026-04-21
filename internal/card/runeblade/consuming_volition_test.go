package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestConsumingVolition_ArcaneDamageNotDealtReturnsBaseAttack: without the arcane-damage clause
// satisfied, the discard rider can't fire regardless of hit likelihood.
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
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard: Red (attack 4) is the only variant
// whose printed attack lands in the likely set ({1,4,7}). With ArcaneDamageDealt set the rider
// fires and Play returns attack+3.
func TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true}
	if got := (ConsumingVolitionRed{}).Play(&s, nil); got != 4+3 {
		t.Errorf("Red with ArcaneDamageDealt: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// TestConsumingVolition_BlockableBaseSuppressesDiscard: Yellow (3) and Blue (2) deliver damage
// multiples the opponent will comfortably block, so the rider doesn't fire even with the
// arcane-damage clause satisfied.
func TestConsumingVolition_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ConsumingVolitionYellow{}, 3},
		{ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{ArcaneDamageDealt: true}
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s with ArcaneDamageDealt: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_RunechantsDontRescue: "When this hits" refers to Consuming Volition's
// own damage reaching the hero. Runechants firing alongside are separate arcane damage and
// don't count toward "this" card hitting, so they can't rescue a blockable variant.
func TestConsumingVolition_RunechantsDontRescue(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true, Runechants: 1}
	if got := (ConsumingVolitionYellow{}).Play(&s, nil); got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
