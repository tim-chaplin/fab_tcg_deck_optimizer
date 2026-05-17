package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// Tests that state.deck and state.hand are reset between permutations so a draw in one
// permutation doesn't poison the next.
func TestPlaySequence_DrawDoesNotPoisonSubsequentPermutations(t *testing.T) {
	top := testutils.RedAttack{}
	deck := []card.Card{top, testutils.BlueAttack{}, testutils.RedAttack{}}
	ctx := NewSequenceContextForTest(hero.Viserai{}, nil, deck, 10, 0, 1)

	// First permutation: Snatch fires, DrawOne pops the top of the deck into Hand.
	_, _, _, _ = ctx.PlaySequence([]card.Card{cards.SnatchRed{}})
	if h := ctx.PermEngine().Hand(); len(h) != 1 || h[0] != top {
		t.Fatalf("after first permutation: Hand = %v, want [top]", h)
	}
	if got := ctx.PermEngine().Deck().Size(); got != len(deck)-1 {
		t.Fatalf("after first permutation: Deck size = %d, want %d (top consumed)",
			got, len(deck)-1)
	}

	// Second permutation: plain attack, no draw. The reset at the top of playSequenceWithMeta
	// must restore state.deck to the original and clear state.hand before this call runs.
	_, _, _, _ = ctx.PlaySequence([]card.Card{testutils.RedAttack{}})
	if h := ctx.PermEngine().Hand(); len(h) != 0 {
		t.Errorf("after second permutation: Hand = %v, want empty (reset lost)", h)
	}
	if got := ctx.PermEngine().Deck().Size(); got != len(deck) {
		t.Errorf("after second permutation: Deck size = %d, want %d (reset lost)",
			got, len(deck))
	}
}

// TestBest_DrawRiderSeesActualDeck pins that the evaluated result depends on the deck
// contents the solver actually receives — two Best calls with identical hands but different
// decks must report distinct end-of-turn State.Hand contents (the cards drawn off the top).
func TestBest_DrawRiderSeesActualDeck(t *testing.T) {
	h := []card.Card{cards.SnatchRed{}}
	deckA := DeckOf(testutils.RedAttack{})
	deckB := DeckOf(testutils.BlueAttack{})

	resA := Best(nil, h, Matchup{IncomingDamage: 0}, deckA, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
	resB := Best(nil, h, Matchup{IncomingDamage: 0}, deckB, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())

	containsID := func(cs []card.Card, id ids.CardID) bool {
		for _, c := range cs {
			if c.ID() == id {
				return true
			}
		}
		return false
	}
	if !containsID(resA.State.Hand(), (testutils.RedAttack{}).ID()) && resA.State.Arsenal() == nil {
		t.Errorf("deck A: drawn RedAttack didn't surface in State.Hand or State.Arsenal: hand=%v arsenal=%v",
			resA.State.Hand(), resA.State.Arsenal())
	}
	if !containsID(resB.State.Hand(), (testutils.BlueAttack{}).ID()) && resB.State.Arsenal() == nil {
		t.Errorf("deck B: drawn BlueAttack didn't surface in State.Hand or State.Arsenal: hand=%v arsenal=%v",
			resB.State.Hand(), resB.State.Arsenal())
	}
}

// Tests an information-leak invariant: with a hand containing a "draw a card" action, the
// chosen hand roles must not depend on Deck[0] — the player commits before seeing the draw.
func TestBest_DeckOrderDoesNotAffectHandRoles(t *testing.T) {
	h := []card.Card{testutils.CostlyDraw{}, testutils.CostlyAttack{}, testutils.PitchOneDR{}}
	deckA := DeckOf(testutils.HugeAttack{}, testutils.PitchOneDR{})
	deckB := DeckOf(testutils.PitchOneDR{}, testutils.HugeAttack{})

	rolesFor := func(summary TurnSummary) map[ids.CardID]deck.Role {
		m := make(map[ids.CardID]deck.Role, len(summary.BestLine))
		for _, a := range summary.BestLine {
			m[a.Card.ID()] = a.Role
		}
		return m
	}

	resA := Best(nil, h, Matchup{IncomingDamage: 0}, deckA, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	resB := Best(nil, h, Matchup{IncomingDamage: 0}, deckB, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())

	rolesA := rolesFor(resA)
	rolesB := rolesFor(resB)
	for id, rA := range rolesA {
		if rB, ok := rolesB[id]; !ok || rA != rB {
			t.Errorf("role differs by deck order for one of the hand cards: deckA role=%s "+
				"deckB role=%s. The solver chose its hand roles based on what it peeked at the "+
				"top of the deck — information the player doesn't have. Lines: A=[%s] B=[%s]",
				rA, rB, FormatBestLine(resA.BestLine), FormatBestLine(resB.BestLine))
			_ = id
		}
	}
}
