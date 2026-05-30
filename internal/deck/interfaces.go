package deck

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

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

// Card is the rich card.Card contract. Deck used to define its own narrow ID+DisplayName
// view to stay decoupled from card, but that decoupling forced []deck.Card<->[]card.Card
// conversions every time the deck handed cards back to the sim — visible as a top entry
// in the from-scratch-anneal alloc profile (PeekTopN's slice copy alone, plus every Opt /
// runOneShuffle / drawNextHand boundary). Aliasing here erases all of them: []deck.Card
// IS []card.Card, so PeekTopN can return a deck-backing view, and Opt's split feeds the
// hero directly with no copy.
type Card = card.Card

// Registry is the legal card / weapon roster Deck constructs against. The two methods
// hand back the full pre-filtered pools (the production registry's NotImplemented /
// Unplayable markers are already excluded), so deck.Random picks from them directly. No
// GetCard or membership predicate on the interface — the production registry isn't
// restructured just to satisfy a lookup deck doesn't need.
type Registry interface {
	LegalCards() []Card
	LegalWeapons() []Weapon
}
