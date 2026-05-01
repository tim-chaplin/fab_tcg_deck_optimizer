package sim

// Forward declarations for the simulator. Each entry below is a package-level function or
// variable whose body / value is supplied by another package at init time so we can break
// would-be import cycles between sim and packages that depend on sim's types.
//
// The pattern is uniform: declare the var here with a sensible default; the providing
// package's init() reassigns it to its real implementation before any caller invokes it.
// Production code always pulls the providing package in transitively (the registry is
// imported by every entry point), so the defaults are only seen in narrowly-scoped tests.
//
// Keeping these in one file makes the cycle-breaking surface easy to audit:
// `grep -r '\bsim\.\(GetCard\|DeckableCards\|AllWeapons\|DisplayName\|ChainStepText\)\b'` lights
// up every consumer.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// DisplayName returns the card name with a pitch-color suffix — "Mauvrion Skies [Y]" for
// a pitch-2 yellow printing. Use anywhere a human-readable identifier needs to
// disambiguate pitch variants (log lines, deck listings, debug printouts).
//
// optimizations.init swaps in a memoised version that hits a per-CardID cache.
var DisplayName = func(c Card) string {
	switch c.Pitch() {
	case 1:
		return c.Name() + " [R]"
	case 2:
		return c.Name() + " [Y]"
	case 3:
		return c.Name() + " [B]"
	}
	return c.Name()
}

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-step
// log line is built from. VERB picks WEAPON ATTACK for weapon-typed cards, ATTACK for
// attack-action cards, DEFENSE REACTION for Defense Reactions, and PLAY for everything
// else; the "from arsenal" suffix tags entries played out of the arsenal slot.
//
// optimizations.init swaps in a memoised version that hits a per-(CardID, FromArsenal)
// cache.
var ChainStepText = func(self *CardState) string {
	types := self.Card.Types()
	var verb string
	switch {
	case types.Has(card.TypeWeapon):
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return DisplayName(self.Card) + ": " + verb
}

// GetCard returns the registered card.Card for id. Populated by registry.init; default
// panics so an early caller (test that forgot to import the registry) gets a clear error
// instead of a silent zero-value Card.
var GetCard = func(id ids.CardID) Card {
	panic("sim.GetCard: registry not loaded — import internal/registry blank to populate the hook")
}

// DeckableCards returns every CardID legal to put in a real deck. Populated by
// registry.init.
var DeckableCards = func() []ids.CardID {
	panic("sim.DeckableCards: registry not loaded — import internal/registry blank to populate the hook")
}

// AllWeapons is the registered weapon roster. Populated by registry.init.
var AllWeapons []Weapon
