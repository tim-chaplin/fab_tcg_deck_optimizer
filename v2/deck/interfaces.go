package deck

import "github.com/tim-chaplin/fab-deck-optimizer/v2/ids"

// Hero is a deck's hero slot. Deck stores Hero through verbatim and never calls into it;
// the simulator's own Hero contract carries every behaviour. Keeping the deck-side surface
// empty means a test can put any value in d.Hero.
type Hero any

// Weapon is what Deck needs from a weapon to validate loadouts and enumerate legal equip
// combinations.
type Weapon interface {
	Name() string
	Hands() int
}

// Card is what Deck needs from a card to enforce per-printing copy budgets and dedupe
// against the user-managed sideboard. DisplayName is the canonical human-readable
// identifier including pitch suffix ("Read the Runes [R]") — sideboard / equipment merges
// compare by it so the three pitch printings of one card stay distinct.
//
// Anything beyond ID and DisplayName (Play, Cost, Attack, …) belongs on the simulator's
// own richer Card contract; consumers that need rich behaviour assert at the read site.
type Card interface {
	ID() ids.CardID
	DisplayName() string
}

// Registry is the legal card / weapon roster Deck constructs against. The two methods
// hand back the full pre-filtered pools (the production registry's NotImplemented /
// Unplayable markers are already excluded), so deck.Random picks from them directly. No
// GetCard or membership predicate on the interface — the production registry isn't
// restructured just to satisfy a lookup deck doesn't need.
type Registry interface {
	LegalCards() []Card
	LegalWeapons() []Weapon
}
