package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized
// for the given chain length. runechantCarryover is wrapped into a Runechant token aura
// on the leaf state so playSequence reads the count off the live aura set, matching
// production.
func newSequenceContextForTest(h hero.Hero, pitched, deckCards []card.Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	d := deck.New(nil, nil, dc)
	leafState := gameengine.GameStateBuilder().
		SetHero(h).
		SetDeck(d).
		Build()
	if runechantCarryover > 0 {
		leafState.CreateAura(aura.NewRunechant(runechantCarryover))
	}
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               d,
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		leafState:          leafState,
	}
}
