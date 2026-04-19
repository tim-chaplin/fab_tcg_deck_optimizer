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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard: Red (attack 4) is the only variant
// whose printed attack lands in the likely set ({1,4,7}). With ArcaneDamageDealt set the rider
// fires and Play returns attack+3.
func TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true}
	if got := (ConsumingVolitionRed{}).Play(&s); got != 4+3 {
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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with ArcaneDamageDealt: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestConsumingVolition_RunechantRescuesBlockableVariants: a single Runechant firing alongside
// is likely to slip through (opponent won't pitch a card to block 1 arcane), so the discard
// rider fires even when the printed attack is blockable.
func TestConsumingVolition_RunechantRescuesBlockableVariants(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true, Runechants: 1}
	if got := (ConsumingVolitionYellow{}).Play(&s); got != 3+3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 6 (runechant slips → rider fires)", got)
	}
}
