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
