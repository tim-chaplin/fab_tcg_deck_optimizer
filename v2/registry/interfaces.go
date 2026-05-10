// Package registry is the master roster of every implemented card, weapon, and hero — the
// deck builder's source of truth for what's legal to put in a deck. See
// docs/dev-standards.md "Registry / sim split" for the package contract and how the sim
// hooks are wired.
package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// Card is the registry-side view of a printed card: identity + display name, the bits
// needed to index, name, and dedupe printings. Name is the printed name without pitch
// suffix ("Aether Slash"); DisplayName includes it ("Aether Slash [R]") so all three
// printings of one card map to distinct entries in the byName index.
type Card interface {
	ID() ids.CardID
	Name() string
	DisplayName() string
}

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
// a new deck or swap one in via mutation. Orthogonal to NotSilverAgeLegal: a card can be
// format-legal yet unimplemented, banned yet fully implemented, or both.
type NotImplemented interface {
	NotImplemented()
}

// Unplayable opts a card or weapon out of the deck-construction pool because the effect
// is so weak the optimizer would never pick it even if fully modelled. Same filtering
// semantics as NotImplemented; the distinction is intent (won't-model vs not-worth-
// modelling). Pre-built hands still play these normally.
type Unplayable interface {
	Unplayable()
}

// NotSilverAgeLegal flags cards / weapons banned in the Silver Age format. The registry
// filters them out of the deck-construction pool. Source of truth is
// data_sources/silver_age_banlist.txt — keep the two in sync.
type NotSilverAgeLegal interface {
	NotSilverAgeLegal()
}
