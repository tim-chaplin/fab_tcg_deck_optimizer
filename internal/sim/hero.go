// Package sim's Hero interface. A FaB hero card — every concrete hero (e.g.
// internal/heroes/Viserai) satisfies this. The methods cover what sim and sim's external
// consumers (cmd/fabsim, internal/deckio, internal/fabrary) need: identity / display, the
// Intelligence hand-draw size, type-set, class, and the per-card-played hook the chain
// runner calls. gameengine has its own narrower Hero interface for what the engine itself
// reads; concrete heroes satisfy both structurally without sim referencing gameengine.Hero.

package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Hero is sim's view of a hero. The method set is a superset of gameengine.Hero — concrete
// heroes (internal/heroes/Viserai, internal/testutils.Hero, the per-test stubs) satisfy
// both interfaces structurally, so sim hands a Hero through to engine.SetHero without
// referencing gameengine.Hero by name.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, g *gameengine.GameEngine, l card.Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}
