package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. Tests use this instead of hand-rolling the context fields so the
// common shape is centralised. Lives in a sim test file (rather than sim_test) so
// exports_test.go's NewSequenceContextForTest wrapper can reach it.
//
// runechantCarryover is wrapped into a priorAuras slice carrying a Runechant token aura,
// matching what production builds from a real previous-turn carryover.
func newSequenceContextForTest(h Hero, pitched, deckCards []card.Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	var priorAuras []Aura
	if runechantCarryover > 0 {
		priorAuras = []Aura{NewRunechantAura(runechantCarryover)}
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
