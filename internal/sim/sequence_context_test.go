package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. runechantCarryover is wrapped into a priorAuras slice carrying a
// Runechant token aura, matching production carryover.
func newSequenceContextForTest(h gameengine.Hero, pitched, deckCards []card.Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	var priorAuras []gameengine.Aura
	if runechantCarryover > 0 {
		priorAuras = []gameengine.Aura{NewRunechantAura(runechantCarryover)}
	}
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               deck.New(nil, nil, dc),
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		priorAuras:         priorAuras,
		carryWinner:        &bufs.carryWinnerScratch,
	}
}
