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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, no aura)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard: Red (attack 4) is the only
// variant whose printed attack lands in the likely set. With AuraCreated set the rider fires.
func TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard(t *testing.T) {
	s := card.TurnState{AuraCreated: true}
	if got := (ReekOfCorruptionRed{}).Play(&s); got != 4+3 {
		t.Errorf("Red with AuraCreated: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// TestReekOfCorruption_AuraPlayedTriggersDiscard: the HasPlayedType(TypeAura) branch satisfies
// the rider the same as AuraCreated.
func TestReekOfCorruption_AuraPlayedTriggersDiscard(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (ReekOfCorruptionRed{}).Play(&s); got != 4+3 {
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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with AuraCreated: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestReekOfCorruption_RunechantRescuesBlockableVariants: a lone Runechant slipping through
// counts as the attack connecting, firing the rider.
func TestReekOfCorruption_RunechantRescuesBlockableVariants(t *testing.T) {
	s := card.TurnState{AuraCreated: true, Runechants: 1}
	if got := (ReekOfCorruptionYellow{}).Play(&s); got != 3+3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 6 (runechant slips → rider fires)", got)
	}
}
