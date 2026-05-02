package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the discard rider stays dormant when ArcaneDamageDealt is false.
func TestConsumingVolition_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ConsumingVolitionRed{}, 4},
		{ConsumingVolitionYellow{}, 3},
		{ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		cs := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, cs)
		testutils.FireOnHitIfLikely(&s, cs)
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that the discard rider fires when ArcaneDamageDealt is set and the attack is likely
// to hit.
func TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard(t *testing.T) {
	s := sim.TurnState{ArcaneDamageDealt: true}
	c := ConsumingVolitionRed{}
	cs := &sim.CardState{Card: c}
	c.Play(&s, cs)
	testutils.FireOnHitIfLikely(&s, cs)
	if got := s.Value; got != 4+3 {
		t.Errorf("Red with ArcaneDamageDealt: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// Tests that the discard rider doesn't fire on blockable variants even with
// ArcaneDamageDealt set.
func TestConsumingVolition_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ConsumingVolitionYellow{}, 3},
		{ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{ArcaneDamageDealt: true}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s with ArcaneDamageDealt: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that co-firing runechants don't rescue a blockable variant — "this hits" reads only
// this card's own damage.
func TestConsumingVolition_RunechantsDontRescue(t *testing.T) {
	s := sim.TurnState{ArcaneDamageDealt: true, Runechants: 1}
	c := ConsumingVolitionYellow{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
