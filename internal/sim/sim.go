// Package sim is a thin compatibility wrapper around internal/deck for
// callers that want a one-shot "shuffle N times and hand me the stats"
// entry point without managing a Deck object.
package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// Run is a convenience wrapper around deck.New(cards).Evaluate(...).
func Run(cards []card.Card, runs int, incomingDamage int, rng *rand.Rand) deck.Stats {
	return deck.New(cards).Evaluate(runs, incomingDamage, rng)
}
