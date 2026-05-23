// Package hero defines the Hero interface a FaB hero card satisfies. Concrete heroes live
// in the sibling internal/hero/heroes subpackage.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Hero is the view of a hero callers need beyond what the engine itself uses. The methods
// after Opt are the trigger-dispatch surface — the same one auras and items expose
// through their embedded trigger.Trigger.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	Opt(cards []card.Card) (top, bottom []card.Card)
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	Matches(types card.TypeSet) bool
	Fire(engine card.GameEngine, logger card.Logger)
}
