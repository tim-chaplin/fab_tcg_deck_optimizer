package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestSutcliffesResearchNotes_EmptyDeck(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (empty deck reveals nothing)", got)
	}
}

func TestSutcliffesResearchNotes_CountsRunebladeAttackActions(t *testing.T) {
	deck := []card.Card{
		testutils.FakeRedAttack().WithTypes(card.TypeRuneblade),
		testutils.FakeRedAction(),
		testutils.FakeRedAttack().WithTypes(card.TypeRuneblade),
	}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := ge.Value(); got != 2 {
		t.Errorf("Red (reveal 3): Play() = %d, want 2 (2 of 3 are Runeblade attack actions)", got)
	}
	if ge.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", ge.RunechantCount())
	}
}

func TestSutcliffesResearchNotes_DeckShorterThanRevealCount(t *testing.T) {
	deck := []card.Card{testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Red (reveal 3, deck 1): Play() = %d, want 1 (only 1 card to reveal)", got)
	}
}

func TestSutcliffesResearchNotes_RunebladeNonAttackIgnored(t *testing.T) {
	// A Runeblade card that isn't an attack action (e.g. Read the Runes: Runeblade + Action, no
	// Attack type) shouldn't count toward the Runechant creation.
	deck := []card.Card{cards.ReadTheRunesRed{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (Runeblade non-attack card shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_NonRunebladeAttackIgnored(t *testing.T) {
	// An attack action that isn't Runeblade-classed shouldn't count.
	deck := []card.Card{testutils.FakeRedAttack()}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SutcliffesResearchNotesRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't count)", got)
	}
}

func TestSutcliffesResearchNotes_VariantRevealCounts(t *testing.T) {
	deck := []card.Card{
		testutils.FakeRedAttack().WithTypes(card.TypeRuneblade),
		testutils.FakeRedAttack().WithTypes(card.TypeRuneblade),
		testutils.FakeRedAttack().WithTypes(card.TypeRuneblade),
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(deck).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
