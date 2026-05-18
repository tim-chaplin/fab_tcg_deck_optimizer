package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
)

// TurnSummary is the sim-runtime working type returned by Best: winning card-role
// assignments plus a *GameState snapshot (hand, deck, arsenal, graveyard, banished, auras,
// items, log entries). Pure data, no rules engine. State's contents depend on the caller —
// production leaves Best's post-chain state; integration-test entry points overwrite with
// a post end-of-turn-cleanup snapshot. deck.BestTurn persists only the durable fields.
type TurnSummary struct {
	BestLine       []deck.CardAssignment
	SwungWeapons   []string
	Value          int
	State          *gameengine.GameState
	IncomingDamage int
	Cacheable      bool
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal.
func (t TurnSummary) ArsenalIn() (deck.CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return deck.CardAssignment{}, false
}

// auraCountByNameInState scans the state's aura list for a token aura with the given
// display name and returns its Count, or zero if no matching entry is present.
func auraCountByNameInState(gs *gameengine.GameState, name string) int {
	if gs == nil {
		return 0
	}
	for _, a := range gs.Auras() {
		if a.CardName() == name {
			return a.Count()
		}
	}
	return 0
}

// itemCountByNameInState is the items counterpart of auraCountByNameInState.
func itemCountByNameInState(gs *gameengine.GameState, name string) int {
	if gs == nil {
		return 0
	}
	for _, i := range gs.Items() {
		if i.CardName() == name {
			return i.Count()
		}
	}
	return 0
}

// auraCountByName scans a sim-concrete aura slice for a token aura by display name.
func auraCountByName(auras []*aura.Aura, name string) int {
	for _, a := range auras {
		if a.CardName() == name {
			return a.Count()
		}
	}
	return 0
}

// itemCountByName scans a sim-concrete item slice for a token item by display name.
func itemCountByName(items []*item.Item, name string) int {
	for _, i := range items {
		if i.CardName() == name {
			return i.Count()
		}
	}
	return 0
}
