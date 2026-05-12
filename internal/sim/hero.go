// Package sim's Hero interface. A FaB hero card — every concrete hero (e.g.
// internal/heroes/Viserai) satisfies this. The methods cover what sim and sim's external
// consumers (cmd/fabsim, internal/deckio, internal/fabrary) need: identity / display, the
// Intelligence hand-draw size, type-set, class, and the per-card-played hook the chain
// runner calls. gameengine has its own narrower Hero interface for what the engine itself
// reads; concrete heroes satisfy both structurally.

package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Hero is sim's view of a hero. Embeds gameengine.Hero so anything satisfying sim.Hero
// also satisfies gameengine.Hero (the engine-level view), letting sim hand a hero through
// to engine.SetHero without a typed conversion.
type Hero interface {
	gameengine.Hero
	ID() ids.HeroID
	Name() string
	Intelligence() int
	OnCardPlayed(played card.Card, g *gameengine.GameEngine, l card.Logger) int
}
