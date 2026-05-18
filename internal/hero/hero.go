// Package hero defines the Hero interface a FaB hero card satisfies. Concrete heroes live
// in the sibling internal/hero/heroes subpackage.
//
// OnCardPlayed takes package-local GameEngine / Logger interfaces (interfaces.go) to avoid
// depending on internal/card.GameEngine. Card-typed arguments and identification constants
// (CardType / TypeSet / Card) are imported from internal/card.
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
