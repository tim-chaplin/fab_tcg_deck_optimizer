package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestSutcliffesResearchNotes_EmptyDeck(t *testing.T) {
	s := &sim.TurnState{}
	(SutcliffesResearchNotesRed{}).Play(s, &sim.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck reveals nothing)", got)
	}
}

func TestSutcliffesResearchNotes_CountsRunebladeAttackActions(t *testing.T) {
	deck := []sim.Card{
		testutils.RunebladeAttack{},
		testutils.NonAttack{},
		testutils.RunebladeAttack{},
	}
	s := &sim.TurnState{Deck: deck}
	(SutcliffesResearchNotesRed{}).Play(s, &sim.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 2 {
		t.Errorf("Red (reveal 3): Play() = %d, want 2 (2 of 3 are Runeblade attack actions)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}

func TestSutcliffesResearchNotes_DeckShorterThanRevealCount(t *testing.T) {
	deck := []sim.Card{testutils.RunebladeAttack{}}
	s := &sim.TurnState{Deck: deck}
	(SutcliffesResearchNotesRed{}).Play(s, &sim.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Red (reveal 3, deck 1): Play() = %d, want 1 (only 1 card to reveal)", got)
	}
}

func TestSutcliffesResearchNotes_RunebladeNonAttackIgnored(t *testing.T) {
	// A Runeblade card that isn't an attack action (e.g. Read the Runes: Runeblade + Action, no
	// Attack type) shouldn't count toward the Runechant creation.
	deck := []sim.Card{ReadTheRunesRed{}}
	s := &sim.TurnState{Deck: deck}
	(SutcliffesResearchNotesRed{}).Play(s, &sim.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (Runeblade non-attack card shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_NonRunebladeAttackIgnored(t *testing.T) {
	// An attack action that isn't Runeblade-classed shouldn't count.
	deck := []sim.Card{testutils.NonRunebladeAttack{}}
	s := &sim.TurnState{Deck: deck}
	(SutcliffesResearchNotesRed{}).Play(s, &sim.CardState{Card: SutcliffesResearchNotesRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_VariantRevealCounts(t *testing.T) {
	deck := []sim.Card{
		testutils.RunebladeAttack{},
		testutils.RunebladeAttack{},
		testutils.RunebladeAttack{},
	}
	cases := []struct {
		c    sim.Card
		want int
	}{
		{SutcliffesResearchNotesRed{}, 3},
		{SutcliffesResearchNotesYellow{}, 2},
		{SutcliffesResearchNotesBlue{}, 1},
	}
	for _, tc := range cases {
		s := &sim.TurnState{Deck: deck}
		tc.c.Play(s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
