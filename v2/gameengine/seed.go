package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// PermutationSeed is the per-permutation source data sim hands to GameEngine.Reset. The
// engine reuses its own internal slice backings across Reset calls, so callers don't pass
// backing arrays — they just supply source slices and the engine appends in.
type PermutationSeed struct {
	Hand                 []card.Card
	Deck                 *deck.Deck // engine copies into its own scratch deck
	Arsenal              card.Card
	Graveyard            []card.Card
	Banished             []card.Card
	OpponentMarked       bool
	Auras                []Aura
	Items                []Item
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

// PersistentSnapshot is GameEngine.Snapshot's return: the persistent state that carries
// into the next turn. Sim's CarryState copies fields out of this snapshot to record the
// winning permutation's end-of-chain state.
type PersistentSnapshot struct {
	Hand           []card.Card
	Deck           *deck.Deck // a fresh Copy() so subsequent permutations don't mutate it
	Arsenal        card.Card
	Graveyard      []card.Card
	Banished       []card.Card
	Auras          []Aura
	Items          []Item
	CardsDrawn     int
	OpponentMarked bool
	LogEntries     []turnlogger.LogEntry
}
