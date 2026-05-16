// Package hero defines the Hero interface a FaB hero card satisfies. Sim and sim's external
// consumers (cmd/fabsim, internal/deckio, internal/fabrary) call methods on hero.Hero
// values; concrete heroes live in this package.
//
// The OnCardPlayed signature uses package-local GameEngine / Logger interfaces (declared
// in interfaces.go) — narrow surfaces that *gameengine.GameEngine and *turnlogger.TurnLogger
// satisfy structurally — so v2/hero doesn't depend on v2/card.GameEngine or v2/card.Logger.
// The card-typed arguments and constants (CardType / TypeSet / Card) are foundational
// identification types and remain imported from v2/card.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Hero is the view of a hero callers need beyond what the engine itself uses.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, ge GameEngine, l Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}
