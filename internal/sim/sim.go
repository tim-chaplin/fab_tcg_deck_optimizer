// Package sim is a thin compatibility wrapper around internal/deck for callers that want a one-shot
// "shuffle N times and hand me the stats" entry point without managing a Deck object.
package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// Run is a convenience wrapper that evaluates the given cards under a default Viserai hero. Useful
// for tests that don't care about hero selection yet.
func Run(cards []card.Card, runs int, incomingDamage int, rng *rand.Rand) deck.Stats {
	return deck.New(hero.Viserai{}, nil, cards).Evaluate(runs, incomingDamage, rng)
}
