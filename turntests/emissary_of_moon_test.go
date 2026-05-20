package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the hand cycle pops the spare hand card to the deck bottom.
func TestEmissaryOfMoon_CyclesHandToDeckBottom(t *testing.T) {
	// BluePitch's ID sorts after Emissary's, so the chain runner's PopHandAt(0) inside
	// the Play hook pops BluePitch. Five-card filler deck keeps BluePitch at the bottom
	// past Emissary's mid-turn DrawOne (1) plus end-of-turn refill (4 more).
	spare := testutils.BluePitch{}
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	})
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.EmissaryOfMoonRed{}, spare})
	if got := summary.State.Deck().NameCounts()[spare.DisplayName()]; got != 1 {
		t.Errorf("deck contains %d copies of the cycled %s, want 1", got, spare.DisplayName())
	}
}
