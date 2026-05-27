package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static
// profile plus a Play hook for on-play behaviour. Implementations live in internal/cards.
type Card interface {
	// ID returns the card's canonical registry identifier, stable within a build.
	ID() ids.CardID
	// Name returns the card's printed name without any pitch-color suffix — all three
	// printings of "Aether Slash" return the same string. Backs synergies and "if you have
	// played a card named X this turn" effects. For display, use DisplayName.
	Name() string
	// DisplayName returns the human-readable identifier including the pitch suffix —
	// "Aether Slash [R]" / "[Y]" / "[B]". Used in log lines, deck listings, debug.
	DisplayName() string
	// Cost returns the card's PRINTED resource cost. This is a static value — the number
	// in the cost gem on the card. Cards whose actual paid cost varies with game state
	// (discounts, runechant-consumes, etc.) implement VariableCost.EffectiveCost(g) and
	// expose the live value via CardState.EffectiveCost(g); Cost() itself is the printed
	// upper bound. Predicates that read "a card with cost N" (Flock's reveal, Demolition
	// Crew's "cost 2 or greater") read Cost() — printed cost is cache-key-safe.
	Cost() int
	Pitch() int
	// Attack is the printed attack value. Conditional bonuses belong in Play.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a TypeSet bitfield. Universal cards
	// fold the active hero's class into the result via g.CurrentHeroClass(); non-Universal
	// cards ignore g.
	Types(g GameEngine) TypeSet
	// GoAgain reports whether playing this card grants an additional action point. Cards
	// printed with "Go again" return true. Hero-conditional cards (Life for a Life, Blow for
	// a Blow, Scar for a Scar) read g.HeroWantsLowerHealth() and return true only when the
	// active hero opts in.
	GoAgain(g GameEngine) bool
	// Play is called when the card resolves — as an attack or as a defense reaction. Cards
	// own card-specific behaviour only: conditional self-buffs, OnHit registration, riders,
	// mid-attack-turn effects. The canonical attack-step credit + log happens in
	// sim.ResolveAttackStep after Play returns — vanilla attacks / DRs have an empty body.
	Play(g GameEngine, l Logger, self *CardState)
}
