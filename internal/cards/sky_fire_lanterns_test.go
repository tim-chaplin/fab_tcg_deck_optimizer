package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func TestSkyFireLanterns_EmptyDeck(t *testing.T) {
	s := &sim.TurnState{}
	(SkyFireLanternsRed{}).Play(s, &sim.CardState{Card: SkyFireLanternsRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck)", got)
	}
}

func TestSkyFireLanterns_MatchingTopCard(t *testing.T) {
	// Red variant (pitch 1) matches a top card with pitch 1.
	s := sim.NewTurnState([]sim.Card{HocusPocusRed{}}, nil)
	(SkyFireLanternsRed{}).Play(s, &sim.CardState{Card: SkyFireLanternsRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Red with Red top: Play() = %d, want 1 (pitch match → create Runechant)", got)
	}
	if s.Runechants != 1 {
		t.Errorf("Runechants = %d, want 1", s.Runechants)
	}
}

func TestSkyFireLanterns_MismatchingTopCard(t *testing.T) {
	// Red variant (pitch 1) doesn't match a Blue top card (pitch 3).
	s := sim.NewTurnState([]sim.Card{HocusPocusBlue{}}, nil)
	(SkyFireLanternsRed{}).Play(s, &sim.CardState{Card: SkyFireLanternsRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Red with Blue top: Play() = %d, want 0 (pitch mismatch)", got)
	}
}

func TestSkyFireLanterns_AllVariantsMatchOwnColor(t *testing.T) {
	cases := []struct {
		lantern sim.Card
		top     sim.Card
	}{
		{SkyFireLanternsRed{}, HocusPocusRed{}},
		{SkyFireLanternsYellow{}, HocusPocusYellow{}},
		{SkyFireLanternsBlue{}, HocusPocusBlue{}},
	}
	for _, tc := range cases {
		s := sim.NewTurnState([]sim.Card{tc.top}, nil)
		tc.lantern.Play(s, &sim.CardState{Card: tc.lantern})
		if got := s.Value; got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (same-color top card)", tc.lantern.Name(), got)
		}
	}
}
