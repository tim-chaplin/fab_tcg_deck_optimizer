// Package hero defines the Hero interface a FaB hero card satisfies. Concrete heroes live
// in the sibling internal/hero/heroes subpackage.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
)

// Hero is the view of a hero callers need beyond what the engine itself uses. The
// embedded trigger.Hook is the dispatch surface FireTriggers walks — the same surface
// auras and items expose through their embedded trigger.Trigger.
type Hero interface {
	trigger.Hook
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	Opt(cards []card.Card) (top, bottom []card.Card)
}
