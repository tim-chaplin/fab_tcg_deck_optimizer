package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Prior is sim's view of the cross-turn carryover state — what the engine starts a turn
// with. Best / EvalOneTurnForTesting / runBestForTurn take a Prior so the chain runner
// can mention sim concrete types (*Aura / *Item / Hero) instead of crossing through
// gameengine.Spec for aura / item state. Internal to chain runner, sim builds a
// gameengine.Spec from Prior's scalar / card-typed fields and seeds auras / items via
// per-entry CreateAura / CreateItem calls.
type Prior struct {
	Hero           Hero
	Arsenal        card.Card
	Banished       []card.Card
	Graveyard      []card.Card
	OpponentMarked bool
	Auras          []*Aura
	Items          []*Item
}
