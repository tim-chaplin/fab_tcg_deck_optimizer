package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// PermutationSeed is the per-permutation source data sim hands to GameEngine.Reset.
//
// Aura / Trigger / Item lists aren't on the seed — Reset clears those lists, and the
// caller follows up with CreateAura / CreateTrigger / CreateItem per entry. That shape
// lets each entry be passed as its concrete struct without slice-conversion (a *sim.Aura
// satisfies Aura, but []*sim.Aura is not assignable to []Aura).
type PermutationSeed struct {
	Hand                 []card.Card
	Deck                 *deck.Deck // engine copies into its own scratch deck
	Arsenal              card.Card
	Graveyard            []card.Card
	Banished             []card.Card
	OpponentMarked       bool
	Pitched              []card.Card
	Defenders            []card.Card
	IncomingDamage       int
	ArcaneIncomingDamage int
	BlockTotal           int
	// Logger is the per-turn log sink: a recording *TurnLogger when materialising the
	// printout, nil during the find-best pass so every Log helper short-circuits via
	// turnlogger's nil-receiver guards.
	Logger *turnlogger.TurnLogger
}
