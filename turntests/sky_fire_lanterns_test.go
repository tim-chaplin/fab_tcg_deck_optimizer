package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestSkyFireLanterns_EmptyDeck(t *testing.T) {
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SkyFireLanternsRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck)", got)
	}
}

func TestSkyFireLanterns_MatchingTopCard(t *testing.T) {
	// Red variant (pitch 1) matches a top card with pitch 1.
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.HocusPocusRed{}}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SkyFireLanternsRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Red with Red top: Play() = %d, want 1 (pitch match → create Runechant)", got)
	}
	if s.RunechantCount() != 1 {
		t.Errorf("Runechants = %d, want 1", s.RunechantCount())
	}
}

func TestSkyFireLanterns_MismatchingTopCard(t *testing.T) {
	// Red variant (pitch 1) doesn't match a Blue top card (pitch 3).
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.HocusPocusBlue{}}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SkyFireLanternsRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Red with Blue top: Play() = %d, want 0 (pitch mismatch)", got)
	}
}

func TestSkyFireLanterns_AllVariantsMatchOwnColor(t *testing.T) {
	cases := []struct {
		lantern card.Card
		top     card.Card
	}{
		{cards.SkyFireLanternsRed{}, cards.HocusPocusRed{}},
		{cards.SkyFireLanternsYellow{}, cards.HocusPocusYellow{}},
		{cards.SkyFireLanternsBlue{}, cards.HocusPocusBlue{}},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{tc.top}).Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.lantern})
		if got := s.Value(); got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (same-color top card)", tc.lantern.Name(), got)
		}
	}
}
