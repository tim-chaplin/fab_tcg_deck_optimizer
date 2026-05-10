package deck

import "github.com/tim-chaplin/fab-deck-optimizer/v2/contracts"

// Hero is a deck's hero slot. Deck stores Hero through verbatim and never calls into it;
// the simulator's own Hero contract carries every behaviour. Keeping the deck-side surface
// empty means a test can put any value in d.Hero.
type Hero any

// Card and Weapon are the shared minimal contracts; consumers needing richer behaviour
// (Play, Cost, Attack, Ability, …) assert to sim.Card / sim.Weapon at the read site.
type Card = contracts.Card
type Weapon = contracts.Weapon

// Registry is the legal card / weapon roster Deck constructs against. The two methods
// hand back the full pre-filtered pools (the production registry's NotImplemented /
// Unplayable markers are already excluded), so deck.Random picks from them directly. No
// GetCard or membership predicate on the interface — the production registry isn't
// restructured just to satisfy a lookup deck doesn't need.
type Registry interface {
	LegalCards() []Card
	LegalWeapons() []Weapon
}
