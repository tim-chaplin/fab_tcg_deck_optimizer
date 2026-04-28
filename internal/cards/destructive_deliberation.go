// Destructive Deliberation — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue
// 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var destructiveDeliberationTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DestructiveDeliberationRed struct{}

func (DestructiveDeliberationRed) ID() ids.CardID          { return ids.DestructiveDeliberationRed }
func (DestructiveDeliberationRed) Name() string            { return "Destructive Deliberation" }
func (DestructiveDeliberationRed) Cost(*sim.TurnState) int { return 2 }
func (DestructiveDeliberationRed) Pitch() int              { return 1 }
func (DestructiveDeliberationRed) Attack() int             { return 5 }
func (DestructiveDeliberationRed) Defense() int            { return 2 }
func (DestructiveDeliberationRed) Types() card.TypeSet     { return destructiveDeliberationTypes }
func (DestructiveDeliberationRed) GoAgain() bool           { return false }

// not implemented: ponder tokens
func (DestructiveDeliberationRed) NotImplemented() {}
func (DestructiveDeliberationRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type DestructiveDeliberationYellow struct{}

func (DestructiveDeliberationYellow) ID() ids.CardID          { return ids.DestructiveDeliberationYellow }
func (DestructiveDeliberationYellow) Name() string            { return "Destructive Deliberation" }
func (DestructiveDeliberationYellow) Cost(*sim.TurnState) int { return 2 }
func (DestructiveDeliberationYellow) Pitch() int              { return 2 }
func (DestructiveDeliberationYellow) Attack() int             { return 4 }
func (DestructiveDeliberationYellow) Defense() int            { return 2 }
func (DestructiveDeliberationYellow) Types() card.TypeSet     { return destructiveDeliberationTypes }
func (DestructiveDeliberationYellow) GoAgain() bool           { return false }

// not implemented: ponder tokens
func (DestructiveDeliberationYellow) NotImplemented() {}
func (DestructiveDeliberationYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type DestructiveDeliberationBlue struct{}

func (DestructiveDeliberationBlue) ID() ids.CardID          { return ids.DestructiveDeliberationBlue }
func (DestructiveDeliberationBlue) Name() string            { return "Destructive Deliberation" }
func (DestructiveDeliberationBlue) Cost(*sim.TurnState) int { return 2 }
func (DestructiveDeliberationBlue) Pitch() int              { return 3 }
func (DestructiveDeliberationBlue) Attack() int             { return 3 }
func (DestructiveDeliberationBlue) Defense() int            { return 2 }
func (DestructiveDeliberationBlue) Types() card.TypeSet     { return destructiveDeliberationTypes }
func (DestructiveDeliberationBlue) GoAgain() bool           { return false }

// not implemented: ponder tokens
func (DestructiveDeliberationBlue) NotImplemented() {}
func (DestructiveDeliberationBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
