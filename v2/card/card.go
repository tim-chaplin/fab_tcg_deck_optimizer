package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
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
	// Cost returns the card's current resource cost given the turn state. Static-cost cards
	// ignore g; cards that read g additionally implement VariableCost so the solver can
	// pre-screen with MinCost / MaxCost bounds.
	Cost(g GameEngine) int
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
	// mid-chain effects. The canonical chain-step credit + log happens in
	// sim.ResolveChainStep after Play returns — vanilla attacks / DRs have an empty body.
	Play(g GameEngine, l Logger, self *CardState)
}
