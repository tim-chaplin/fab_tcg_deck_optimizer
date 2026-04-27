package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestSutcliffesResearchNotes_EmptyDeck(t *testing.T) {
	s := &card.TurnState{}
	(SutcliffesResearchNotesRed{}).Play(s, &card.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck reveals nothing)", got)
	}
}

func TestSutcliffesResearchNotes_CountsRunebladeAttackActions(t *testing.T) {
	deck := []card.Card{
		stubRunebladeAttack{},
		stubNonAttack{},
		stubRunebladeAttack{},
	}
	s := card.NewTurnState(deck, nil)
	(SutcliffesResearchNotesRed{}).Play(s, &card.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 2 {
		t.Errorf("Red (reveal 3): Play() = %d, want 2 (2 of 3 are Runeblade attack actions)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}

func TestSutcliffesResearchNotes_DeckShorterThanRevealCount(t *testing.T) {
	deck := []card.Card{stubRunebladeAttack{}}
	s := card.NewTurnState(deck, nil)
	(SutcliffesResearchNotesRed{}).Play(s, &card.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Red (reveal 3, deck 1): Play() = %d, want 1 (only 1 card to reveal)", got)
	}
}

func TestSutcliffesResearchNotes_RunebladeNonAttackIgnored(t *testing.T) {
	// A Runeblade card that isn't an attack action (e.g. Read the Runes: Runeblade + Action, no
	// Attack type) shouldn't count toward the Runechant creation.
	deck := []card.Card{ReadTheRunesRed{}}
	s := card.NewTurnState(deck, nil)
	(SutcliffesResearchNotesRed{}).Play(s, &card.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (Runeblade non-attack card shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_NonRunebladeAttackIgnored(t *testing.T) {
	// An attack action that isn't Runeblade-classed shouldn't count.
	deck := []card.Card{stubNonRunebladeAttack{}}
	s := card.NewTurnState(deck, nil)
	(SutcliffesResearchNotesRed{}).Play(s, &card.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_VariantRevealCounts(t *testing.T) {
	deck := []card.Card{
		stubRunebladeAttack{},
		stubRunebladeAttack{},
		stubRunebladeAttack{},
	}
	cases := []struct {
		c    card.Card
		want int
	}{
		{SutcliffesResearchNotesRed{}, 3},
		{SutcliffesResearchNotesYellow{}, 2},
		{SutcliffesResearchNotesBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.NewTurnState(deck, nil)
		tc.c.Play(s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
