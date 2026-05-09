package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Deck is a candidate deck the simulator evaluates. Aliased to v2/deck.Deck so existing
// sim consumers continue to compile against *sim.Deck while the canonical type lives in
// the deck package. Cards / Weapons hold deck.Card / deck.Weapon — sim's hot paths assert
// to the richer sim.Card / sim.Weapon at the read site.
type Deck = deck.Deck

// NotImplementedReplacement aliases v2/deck.NotImplementedReplacement so the
// SanitizeNotImplemented return type stays accessible under the sim namespace for the
// handful of consumers that surface it (fabsim's sanitize warning, deckio tests).
type NotImplementedReplacement = deck.NotImplementedReplacement

// New constructs a Deck from sim's richer Hero / Weapon / Card interfaces. Mutation
// generators and tests build cards as sim values; this wrapper widens them to deck.* slots
// without callers hand-rolling the per-element copy.
func New(h Hero, weapons []Weapon, cards []Card) *Deck {
	dw := make([]deck.Weapon, len(weapons))
	for i, w := range weapons {
		dw[i] = w
	}
	dc := make([]deck.Card, len(cards))
	for i, c := range cards {
		dc[i] = c
	}
	return deck.New(h, dw, dc)
}

// Random constructs a random legal deck for h against the production registry. Wraps
// deck.Random with sim's registry hooks so callers don't have to import the registry
// directly (would cycle through cards → sim → registry).
func Random(h Hero, size, maxCopies int, rng *rand.Rand, legal func(Card) bool) *Deck {
	var legalDeck func(deck.Card) bool
	if legal != nil {
		legalDeck = func(c deck.Card) bool { return legal(c.(Card)) }
	}
	return deck.Random(h, size, maxCopies, rng, legalDeck, simRegistry{})
}

// simRegistry satisfies deck.Registry by routing through sim's forward-declared
// RegistryLegal* hooks so sim-side callers don't import the registry directly.
type simRegistry struct{}

func (simRegistry) LegalCards() []deck.Card {
	src := RegistryLegalCards()
	out := make([]deck.Card, len(src))
	for i, c := range src {
		out[i] = c
	}
	return out
}

func (simRegistry) LegalWeapons() []deck.Weapon { return RegistryLegalWeapons() }

// legalPool returns the registered card IDs the deck builder can pick from, optionally
// filtered by legal. The marker-exclusion (NotImplemented / Unplayable / NotSilverAgeLegal)
// happens in the registry; this wrapper just maps the registry's Card slice down to IDs
// for sim's mutation enumerator.
func legalPool(legal func(Card) bool) []ids.CardID {
	pool := RegistryLegalCards()
	out := make([]ids.CardID, 0, len(pool))
	for _, c := range pool {
		if legal != nil && !legal(c) {
			continue
		}
		out = append(out, c.ID())
	}
	return out
}

// RegistryLegalCards / RegistryLegalWeapons are forward-declared (capitalised so the
// providing package can rebind them — see forward_declarations.go) to break the would-be
// sim → registry import cycle. registry.init points them at its own Registry methods.
var RegistryLegalCards = func() []Card {
	panic("sim.RegistryLegalCards: registry not loaded — import internal/registry blank to populate the hook")
}

var RegistryLegalWeapons = func() []deck.Weapon {
	panic("sim.RegistryLegalWeapons: registry not loaded — import internal/registry blank to populate the hook")
}
