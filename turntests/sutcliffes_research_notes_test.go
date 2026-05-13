package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestSutcliffesResearchNotes_EmptyDeck(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck reveals nothing)", got)
	}
}

func TestSutcliffesResearchNotes_CountsRunebladeAttackActions(t *testing.T) {
	deck := []card.Card{
		testutils.RunebladeAttack{},
		testutils.NonAttack{},
		testutils.RunebladeAttack{},
	}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := s.Value(); got != 2 {
		t.Errorf("Red (reveal 3): Play() = %d, want 2 (2 of 3 are Runeblade attack actions)", got)
	}
	if s.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", s.RunechantCount())
	}
}

func TestSutcliffesResearchNotes_DeckShorterThanRevealCount(t *testing.T) {
	deck := []card.Card{testutils.RunebladeAttack{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Red (reveal 3, deck 1): Play() = %d, want 1 (only 1 card to reveal)", got)
	}
}

func TestSutcliffesResearchNotes_RunebladeNonAttackIgnored(t *testing.T) {
	// A Runeblade card that isn't an attack action (e.g. Read the Runes: Runeblade + Action, no
	// Attack type) shouldn't count toward the Runechant creation.
	deck := []card.Card{cards.ReadTheRunesRed{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (Runeblade non-attack card shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_NonRunebladeAttackIgnored(t *testing.T) {
	// An attack action that isn't Runeblade-classed shouldn't count.
	deck := []card.Card{testutils.NonRunebladeAttack{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_VariantRevealCounts(t *testing.T) {
	deck := []card.Card{
		testutils.RunebladeAttack{},
		testutils.RunebladeAttack{},
		testutils.RunebladeAttack{},
	}
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.SutcliffesResearchNotesRed{}, 3},
		{cards.SutcliffesResearchNotesYellow{}, 2},
		{cards.SutcliffesResearchNotesBlue{}, 1},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
