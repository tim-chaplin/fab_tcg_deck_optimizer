// Package hero defines the Hero interface a FaB hero card satisfies. Sim and its external
// consumers (cmd/fabsim, internal/textio) call methods on hero.Hero values; concrete heroes
// live in the sibling internal/hero/heroes subpackage.
//
// The OnCardPlayed signature uses package-local GameEngine / Logger interfaces (declared
// in interfaces.go) so internal/hero doesn't depend on internal/card.GameEngine or
// internal/card.Logger. The card-typed arguments and constants (CardType / TypeSet / Card)
// are foundational identification types and remain imported from internal/card.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
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
