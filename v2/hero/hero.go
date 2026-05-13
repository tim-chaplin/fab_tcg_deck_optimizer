// Package hero defines the Hero interface a FaB hero card satisfies. Sim and sim's external
// consumers (cmd/fabsim, internal/deckio, internal/fabrary) call methods on hero.Hero
// values; concrete heroes live in internal/heroes. gameengine has its own narrower Hero
// interface; concrete heroes satisfy both structurally so callers can hand a hero.Hero
// through to engine.SetHero without going through gameengine.Hero by name.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Hero is the view of a hero callers need beyond what the engine itself uses. OnCardPlayed
// takes a card.GameEngine — concrete heroes don't import gameengine because *GameEngine
// satisfies card.GameEngine.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, ge card.GameEngine, l card.Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}
