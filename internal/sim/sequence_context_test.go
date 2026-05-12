package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. runechantCarryover is wrapped into a Runechant token aura on the
// leaf engine so playSequence reads the count off the live aura set, matching production.
func newSequenceContextForTest(h Hero, pitched, deckCards []card.Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	d := deck.New(nil, nil, dc)
	leafEngine := gameengine.NewFromSpec(gameengine.Spec{Hero: h})
	leafEngine.SetDeck(d)
	if runechantCarryover > 0 {
		leafEngine.CreateAura(NewRunechantAura(runechantCarryover))
	}
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               d,
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		leafEngine:         leafEngine,
	}
}
