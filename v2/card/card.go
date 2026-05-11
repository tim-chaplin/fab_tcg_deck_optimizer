package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's
// static profile plus a Play hook for on-play behaviour. Implementations live in
// internal/cards (and similar leaf packages); the sim's chain runner consumes them
// through this interface.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build.
	// Lets callers key maps / slices on cards without string-hashing Name().
	ID() ids.CardID
	// Name returns the card's printed name without any pitch-color suffix — all three
	// printings of "Aether Slash" return the same string. Cards comparing by name
	// (synergies, "if you have played a card named X this turn" effects) use this
	// directly. For display, callers route through DisplayName which appends the
	// pitch tag.
	Name() string
	// DisplayName returns the human-readable identifier including the pitch suffix —
	// "Aether Slash [R]", "Aether Slash [Y]", "Aether Slash [B]". Use this anywhere
	// a printout needs to disambiguate pitch printings (log lines, deck listings,
	// debug).
	DisplayName() string
	// Cost returns the card's current resource cost given the turn state. Cards with
	// a static printed cost ignore g and return a constant; cards that read g (e.g.
	// discount-per-token effects) additionally implement VariableCost so the solver
	// can pre-screen with cheap MinCost / MaxCost bounds before enumerating chain
	// permutations.
	Cost(g GameEngine) int
	Pitch() int
	// Attack is the printed attack value. Conditional bonuses belong in Play, not
	// here.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a card.TypeSet bitfield,
	// e.g. card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack).
	Types() card.TypeSet
	// GoAgain reports whether playing this card grants an additional action point.
	// Cards printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card resolves — as an attack or as a defense reaction.
	// Cards own card-specific behaviour only: conditional self-buffs (flip
	// self.BonusAttack / self.BonusDefense), riders (l.AppendPostTrigger sub-lines),
	// OnHit registration, mid-chain effects. The standard "credit EffectiveAttack /
	// capped EffectiveDefense to g and emit the <Card>: <VERB> (+N) chain step"
	// mechanic happens in sim.ResolveChainStep after Play returns — vanilla attack /
	// DR cards have an empty Play body.
	Play(g GameEngine, l Logger, self *CardState)
}
