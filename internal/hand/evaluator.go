package hand

// Entry points for hand evaluation. Best / BestWithTriggers compute the optimal turn line
// for a given hand against an opponent attacking for incomingDamage. The Evaluator type is
// a no-op wrapper kept so concurrent callers can construct per-goroutine instances; every
// call allocates fresh scratch state.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Best returns the optimal TurnSummary for the given hand against an opponent that will
// attack for incomingDamage on their next turn. Equipped weapons may be swung for their Cost
// if resources allow.
//
// Cards partition into five roles: Pitch (resource), Attack (played, may extend chain),
// Defend (blocks plus DR Plays), Held (stays in hand for next turn), Arsenal (moves to or
// stays in the arsenal slot at end of turn). Pitch resources split across attack / defense
// phases since resources don't carry between turns.
//
// arsenalCardIn is the card sitting in the arsenal slot at start of turn (nil if empty).
// runechantCarryover is the Runechant token count carrying in from the previous turn.
// TurnSummary.State.Runechants is the count at end of the chosen chain; feed it back as
// the next turn's carryover.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is Best plus an explicit priorAuraTriggers input — the AuraTriggers
// carrying in from the previous turn. Mid-chain triggers (Malefic Incantation's
// TriggerAttackAction rune, etc.) may fire and contribute damage to this turn's Value.
func BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
}

// Best is the method form of the package-level Best.
func (e *Evaluator) Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return e.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is the method form of the package-level BestWithTriggers.
func (e *Evaluator) BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	return e.bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
}

// Evaluator is a placeholder for per-goroutine state. Currently empty — every call allocates
// fresh scratch — but kept so concurrent callers can construct one if scratch caching needs
// to be reintroduced behind the same surface.
type Evaluator struct{}

// NewEvaluator returns a fresh Evaluator. Safe for concurrent use across goroutines.
func NewEvaluator() *Evaluator { return &Evaluator{} }

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to
// construct their own.
var sharedEvaluator = NewEvaluator()
