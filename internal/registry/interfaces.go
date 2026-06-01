// Package registry is the master roster of every implemented card, weapon, and hero — the
// deck builder's source of truth for what's legal to put in a deck. See
// docs/dev-standards.md "Registry / sim split" for the package contract and how the sim
// hooks are wired.
package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Card is the rich card.Card contract — aliased here for the same reason deck.Card is:
// the original narrow ID/Name/DisplayName view forced []registry.Card<->[]card.Card
// conversions at every sim boundary. Aliasing erases them so deck-construction pools and
// the sim share one card type end-to-end.
type Card = card.Card

// Hero is the registry-side view of a hero: just the printed name, used as the byName
// lookup key.
type Hero interface {
	Name() string
}

// Weapon is the registry-side view of a weapon: enough to index by name and validate
// loadouts.
type Weapon interface {
	Name() string
	Hands() int
}

// NotImplemented opts a card or weapon out of the deck-construction pool because the
// simulator doesn't model its printed effect (Gold / Silver / Copper token economies,
// Landmarks, …). A NotImplemented card is still valid in pre-built hands — it evaluates
// using its static Attack / Pitch / Defense — but the optimizer won't introduce it into
// a new deck or swap one in via mutation. Orthogonal to format legality (the internal/format
// banlist): a card can be format-legal yet unimplemented, banned yet fully implemented, or
// both.
type NotImplemented interface {
	NotImplemented()
}

// Unplayable opts a card or weapon out of the deck-construction pool because the effect
// is so weak the optimizer would never pick it even if fully modelled. Same filtering
// semantics as NotImplemented; the distinction is intent (won't-model vs not-worth-
// modelling). Pre-built hands still play these normally.
//
// Format legality (the Silver Age banlist) is not a marker: it lives in internal/format, and
// LegalCards / LegalWeapons apply it via format.IsCardLegal.
type Unplayable interface {
	Unplayable()
}
