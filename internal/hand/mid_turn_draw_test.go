package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestPlaySequence_DrawDoesNotPoisonSubsequentPermutations pins the per-permutation reset of
// state.Deck and state.Drawn. Two back-to-back playSequence calls share one sequenceContext —
// the first fires a draw-rider card (Snatch), the second plays a plain attack. After the
// second call finishes, state.Deck must be back to the original and state.Drawn empty; if the
// reset weren't wired in, the second permutation would start from an already-consumed deck
// and an inherited Drawn slice.
func TestPlaySequence_DrawDoesNotPoisonSubsequentPermutations(t *testing.T) {
	top := fake.RedAttack{}
	deck := []card.Card{top, fake.BlueAttack{}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, deck, 10, 0, 1)

	// First permutation: Snatch fires, DrawOne consumes the top of the deck.
	_, _, _, _ = ctx.playSequence([]card.Card{generic.SnatchRed{}}, nil, nil)
	if len(ctx.bufs.state.Drawn) != 1 || ctx.bufs.state.Drawn[0] != top {
		t.Fatalf("after first permutation: Drawn = %v, want [top]", ctx.bufs.state.Drawn)
	}
	if len(ctx.bufs.state.Deck) != len(deck)-1 {
		t.Fatalf("after first permutation: Deck len = %d, want %d (top consumed)",
			len(ctx.bufs.state.Deck), len(deck)-1)
	}

	// Second permutation: plain attack, no draw. The reset at the top of playSequenceWithMeta
	// must restore state.Deck to the original and clear state.Drawn before this call runs.
	_, _, _, _ = ctx.playSequence([]card.Card{fake.RedAttack{}}, nil, nil)
	if len(ctx.bufs.state.Drawn) != 0 {
		t.Errorf("after second permutation: Drawn = %v, want empty (reset lost)", ctx.bufs.state.Drawn)
	}
	if len(ctx.bufs.state.Deck) != len(deck) {
		t.Errorf("after second permutation: Deck len = %d, want %d (reset lost)",
			len(ctx.bufs.state.Deck), len(deck))
	}
}

// TestBest_DrawRiderBypassesMemo pins the NoMemo discipline on mid-turn-draw cards: the
// evaluated result must depend on the deck contents, not just the hand, so `memoKey` (which
// doesn't include the deck) can't cache one deck's outcome into another's. Two Best calls
// with identical hands but different decks must report distinct Drawn cards.
func TestBest_DrawRiderBypassesMemo(t *testing.T) {
	h := []card.Card{generic.SnatchRed{}}
	deckA := []card.Card{fake.RedAttack{}}
	deckB := []card.Card{fake.BlueAttack{}}

	resA := Best(hero.Viserai{}, nil, h, 0, deckA, 0, nil)
	resB := Best(hero.Viserai{}, nil, h, 0, deckB, 0, nil)

	if len(resA.Drawn) != 1 || resA.Drawn[0].Card != (fake.RedAttack{}) {
		t.Errorf("deck A: Drawn = %v, want [RedAttack]", resA.Drawn)
	}
	if len(resB.Drawn) != 1 || resB.Drawn[0].Card != (fake.BlueAttack{}) {
		t.Errorf("deck B: Drawn = %v, want [BlueAttack] (memo collision bleeding deck A's result here)", resB.Drawn)
	}
}

// TestBest_DeckOrderDoesNotAffectHandRoles pins an information-leak invariant in the solver.
//
// Problem: when a hand contains a "draw a card" action, the evaluator doesn't know which card
// the player will see on top of the deck until the draw actually fires. In the real game the
// player has to commit to a line first — play the draw hoping for something useful, or don't —
// then live with whatever comes off the top. The current solver, though, evaluates every
// permutation with full visibility into Deck[0] and lets mid-turn-drawn cards be pitched or
// played as chain extensions (bestSequence's extension loop in playSequenceWithMeta). That
// means the best line it picks can genuinely depend on the identity of the top card: with a
// fantastic attack on top it'll play the draw, with a defense reaction on top it'll skip the
// draw and play something else instead. The player, reordering what's in the same deck, would
// see the same choice offered up — the evaluator has effectively cheated by peeking.
//
// The test: fix the hand and flip two deck orderings. The hand roles have to match. The draw
// card is allowed to be played or not; the invariant is that the choice can't flip as a
// function of deck order alone.
func TestBest_DeckOrderDoesNotAffectHandRoles(t *testing.T) {
	h := []card.Card{fake.CostlyDraw{}, fake.CostlyAttack{}, fake.PitchOneDR{}}
	deckA := []card.Card{fake.HugeAttack{}, fake.PitchOneDR{}}
	deckB := []card.Card{fake.PitchOneDR{}, fake.HugeAttack{}}

	rolesFor := func(summary TurnSummary) map[card.ID]Role {
		m := make(map[card.ID]Role, len(summary.BestLine))
		for _, a := range summary.BestLine {
			m[a.Card.ID()] = a.Role
		}
		return m
	}

	resA := Best(stubHero{}, nil, h, 0, deckA, 0, nil)
	resB := Best(stubHero{}, nil, h, 0, deckB, 0, nil)

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
