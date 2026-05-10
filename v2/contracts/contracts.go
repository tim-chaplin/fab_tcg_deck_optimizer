// Package contracts holds the minimal Card / Weapon interfaces shared by v2/deck and
// v2/registry. See docs/dev-standards.md "Registry / sim split".
package contracts

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// Card is the minimal view of a printed card: identity + display name. Name is the
// printed name without pitch suffix ("Aether Slash"); DisplayName includes it
// ("Aether Slash [R]") so the three pitch printings of one card map to distinct names.
type Card interface {
	ID() ids.CardID
	Name() string
	DisplayName() string
}

// Weapon is the minimal view of a weapon: enough to validate loadouts and dedupe by name.
type Weapon interface {
	Name() string
	Hands() int
}
