package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized
// for the given attack-turn length. runechantCarryover is wrapped into a Runechant token aura
// on the leaf state so playSequence reads the count off the live aura set, matching
// production.
func newSequenceContextForTest(h hero.Hero, pitched, deckCards []card.Card, resourceBudget, runechantCarryover, attackTurnLen int) *sequenceContext {
	bufs := newAttackBufs(attackTurnLen, 0, nil, gameengine.NewPrewarmedPool())
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	d := deck.New(nil, nil, dc)
	builder := gameengine.GameStateBuilder().
		SetHero(h).
		SetDeck(d)
	if runechantCarryover > 0 {
		builder = builder.AddAura(token.NewRunechant(runechantCarryover))
	}
	leafState := builder.Build()
	return &sequenceContext{
		pitched:            wrapPitchedCards(pitched),
		deck:               d,
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		leafState:          leafState,
	}
}
