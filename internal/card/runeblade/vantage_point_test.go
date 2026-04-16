package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestVantagePoint_BaseDamageNoAura(t *testing.T) {
	// No aura → just printed power, Overpower stays false.
	cases := []struct {
		c    card.Card
		base int
	}{
		{VantagePointRed{}, 7},
		{VantagePointYellow{}, 6},
		{VantagePointBlue{}, 5},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
		if s.Overpower {
			t.Errorf("%s: Overpower should stay false when no aura", tc.c.Name())
		}
	}
}

func TestVantagePoint_AuraPlayedSetsOverpower(t *testing.T) {
	// Aura in CardsPlayed → Overpower flag set; damage unchanged since Overpower isn't consumed
	// by the solver (incoming damage is a flat opponent profile, not blocked).
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (VantagePointRed{}).Play(&s); got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !s.Overpower {
		t.Errorf("Overpower should be set when an aura was played")
	}
}

func TestVantagePoint_AuraCreatedSetsOverpower(t *testing.T) {
	// AuraCreated flag (e.g. from an earlier Runechant-creating card) also triggers Overpower.
	s := card.TurnState{AuraCreated: true}
	if got := (VantagePointRed{}).Play(&s); got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !s.Overpower {
		t.Errorf("Overpower should be set when AuraCreated is true")
	}
}
