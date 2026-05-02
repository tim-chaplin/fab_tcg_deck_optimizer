package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the discard rider stays dormant without an aura played or created.
func TestReekOfCorruption_NoAuraReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ReekOfCorruptionRed{}, 4},
		{ReekOfCorruptionYellow{}, 3},
		{ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		cs := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, cs)
		testutils.FireOnHitIfLikely(&s, cs)
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, no aura)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that the discard rider fires with AuraCreated set on a likely-hit attack.
func TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard(t *testing.T) {
	s := sim.TurnState{AuraCreated: true}
	c := ReekOfCorruptionRed{}
	cs := &sim.CardState{Card: c}
	c.Play(&s, cs)
	testutils.FireOnHitIfLikely(&s, cs)
	if got := s.Value; got != 4+3 {
		t.Errorf("Red with AuraCreated: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// Tests that an aura earlier in CardsPlayed satisfies the rider precondition.
func TestReekOfCorruption_AuraPlayedTriggersDiscard(t *testing.T) {
	s := sim.TurnState{CardsPlayed: []sim.Card{testutils.Aura{}}}
	c := ReekOfCorruptionRed{}
	cs := &sim.CardState{Card: c}
	c.Play(&s, cs)
	testutils.FireOnHitIfLikely(&s, cs)
	if got := s.Value; got != 4+3 {
		t.Errorf("Play() = %d, want %d (aura earlier in chain triggers rider)", got, 4+3)
	}
}

// Tests that the discard rider doesn't fire on blockable variants even with AuraCreated.
func TestReekOfCorruption_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ReekOfCorruptionYellow{}, 3},
		{ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{AuraCreated: true}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s with AuraCreated: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that co-firing runechants don't rescue a blockable variant — "this hits" reads only
// this card's own damage.
func TestReekOfCorruption_RunechantsDontRescue(t *testing.T) {
	s := sim.TurnState{AuraCreated: true, Runechants: 1}
	c := ReekOfCorruptionYellow{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
