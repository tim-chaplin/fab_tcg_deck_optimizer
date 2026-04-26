package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestReekOfCorruption_NoAuraReturnsBaseAttack: without an aura played or created this turn the
// discard rider can't fire, regardless of hit likelihood.
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
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, no aura)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard: Red (attack 4) is the only
// variant whose printed attack lands in the likely set. With AuraCreated set the rider fires.
func TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard(t *testing.T) {
	s := card.TurnState{AuraCreated: true}
	c := ReekOfCorruptionRed{}
	c.Play(&s, &card.CardState{Card: c})
	if got := s.Value; got != 4+3 {
		t.Errorf("Red with AuraCreated: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// TestReekOfCorruption_AuraPlayedTriggersDiscard: the HasPlayedType(TypeAura) branch satisfies
// the rider the same as AuraCreated.
func TestReekOfCorruption_AuraPlayedTriggersDiscard(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	c := ReekOfCorruptionRed{}
	c.Play(&s, &card.CardState{Card: c})
	if got := s.Value; got != 4+3 {
		t.Errorf("Play() = %d, want %d (aura earlier in chain triggers rider)", got, 4+3)
	}
}

// TestReekOfCorruption_BlockableBaseSuppressesDiscard: Yellow (3) and Blue (2) are blockable
// totals the opponent won't let through, so the on-hit rider doesn't fire.
func TestReekOfCorruption_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ReekOfCorruptionYellow{}, 3},
		{ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{AuraCreated: true}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s with AuraCreated: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_RunechantsDontRescue: "When this hits" is strictly about this card's own
// damage reaching the hero. Runechants firing alongside are separate arcane damage and don't
// count toward "this" card hitting.
func TestReekOfCorruption_RunechantsDontRescue(t *testing.T) {
	s := card.TurnState{AuraCreated: true, Runechants: 1}
	c := ReekOfCorruptionYellow{}
	c.Play(&s, &card.CardState{Card: c})
	if got := s.Value; got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
